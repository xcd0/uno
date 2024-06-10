package core

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/pkg/errors"
)

// サーバーが保持するゲームの情報。参加者には見せない。
type State struct {
	Color          int      // -1の時任意
	Number         int      // -1の時任意
	Draw           int      // 0以下の時任意
	Turn           int      // CW:1 CCW:-1
	NowPlayerIndex int      // 今の人のID
	Players        []string // 参加者
	NPCs           []bool   // NPCかどうか
	Hands          [][]Card // 手札
	DrawnCardID    []int    // 新規に引いたカードのID。draw2等で複数引いている場合がある。
	Deck           []Card   // 山札
	Discard        []Card   // 捨て札
	LastCard       *Card    // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
	SkipCount      int      // 全員手札が出せない謎状況になったら引き分けにする。
	Rule           UnoRule  // 使用中のルール
	Setting        *Setting // サーバーの設定
	UnoCalled      []bool   // 他プレイヤーがUNO!とコールしたかどうかリスト。
}

func NewState(players []string, numberOfPlayers int, rule *UnoRule) *State {
	return (&State{}).Init(players, numberOfPlayers, GetCards(rule), rule)
}

func (state State) Init(players []string, numberOfPlayers int, cards []Card, rule *UnoRule) *State {

	if len(players) == 0 {
		players = []string{"You", "NPC1", "NPC2", "NPC3"}
	}
	npcs := make([]bool, 0, 20)
	if numberOfPlayers != 0 {
		numberOfPlayers = max(len(players), numberOfPlayers)
		for i := 0; i < len(players); i++ {
			npcs = append(npcs, false)
		}
		if numberOfPlayers > len(players) {
			for i := 0; i < numberOfPlayers-len(players); i++ {
				players = append(players, fmt.Sprintf("P%v", len(players)+i))
				npcs = append(npcs, true)
			}
		}
	}

	// 参加者の並びをシャッフルする。
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	state = State{
		Color:          -1,                      // -1の時任意
		Number:         -1,                      // -1の時任意
		Draw:           0,                       // 0以下の時任意
		Turn:           1,                       // CW:1 CCW:-1
		NowPlayerIndex: rand.Intn(len(players)), // 今の人のID

		Players:   players,                      // 参加者
		NPCs:      npcs,                         // NPCかどうか
		Hands:     make([][]Card, len(players)), // 手札
		Deck:      cards,                        // 山札
		Discard:   make([]Card, 0, len(cards)),  // 捨て札
		LastCard:  nil,                          // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
		SkipCount: 0,                            // 全員手札が出せない謎状況になったら引き分けにする。
		Rule:      *rule,                        // 使用中のルール
		Setting:   nil,                          // サーバーの設定
		UnoCalled: make([]bool, len(players)),   // 他プレイヤーがUNO!とコールしたかどうかリスト。
	}

	return &state
}
func (s *State) InitGameState() *State {
	s.ShuffleDeck()
	s.DealHandsOnInit()
	return s
}
func (s *State) ShuffleDeck() *State {
	// shuffle
	ShuffleCards(s.Deck)
	CryptoRandShuffle(s.Deck)
	return s
}
func (s *State) DealHandsOnInit() *State {

	if s.Rule.EnableAlgoriStart { // R.algori.O1 ALGORIルールの開始方法を使用する。
		// - 順番(回る方向)と最初のプレーヤーはランダム。
		// - 最初のカードの決定方法が特殊。
		//   - ディーラーが1枚出す。場札(最初の場のカードを決めるためにディーラーが配ったカード)とする。
		//   - 場札が
		//     - 色が決まらないカード(ワイルドドロー4、シャッフルワイルド、ホワイトワイルド)の場合は、山札のランダムな位置に戻して、再度ディーラーが場札を出す。
		//     - スキップの場合、スキップの効果が適用される。2番目のプレーヤーからスタート。
		//     - ドロー2の時、ドロー2の効果が適用される。1番目のプレーヤーが2枚引き、2番目のプレイヤーからスタート。
		//     - リバースの時、リバースの効果が適用される。3番目のプレーヤーからリバースしてスタート。
		//   - 場札が決まった後に手札を配る。(公式ルールと異なる。)
		// - 1ターンは5秒以内。
		// - 出せるカードの有無に関わらず(手札が25枚未満の時)カードを1枚引くことができるが、引いたときは引いたカードのみ出せる。
		// - 不正操作の場合(時間切れで出した、誤ったカードを出した、色を指定すべき時に指定しなかった、色を指定すべきでないとき指定した、手札が複数枚の時にUNOコールなど)
		//   - 罰則で2枚引く。出したカードは手札に戻る。スキップされる。
		// - 場札と山札両方が無くなった場合(全て手札にある状態。場札は1枚だと思われる。)、10周経過して誰も出さない場合、勝者なしとして全員敗者扱いする。
		// - ワイルドチャレンジが成功した時、ワイルドドロー4は手札に戻される。

		// - 順番(回る方向)と最初のプレーヤーはランダム。
		s.Turn = RandomSign()

		// 最初のカードの決定。ディーラーが1枚出す。場札(最初の場のカードを決めるためにディーラーが配ったカード)とする。
		var card Card
		for {
			if !CheckDeck(s) { // 山札が弾ける状態かチェック
				panic(errors.Errorf("%v", "バグ")) // 謎。バグ。
			}
			var ok bool
			if card, ok = pop(&s.Deck); !ok {
				panic(errors.Errorf("%v", "バグ")) // 謎。バグ。
			}
			if card.IsWild() { // - 場札が色が決まらないカードの場合は、山札のランダムな位置に戻して、再度ディーラーが場札を出す。
				InsertRandomCard(&s.Deck, card)
				continue
			}
			// 場札にする。
			push(&s.Deck, card)
		}
		// - 場札がスキップの場合、スキップの効果が適用される。2番目のプレーヤーからスタート。
		if card.IsSkip() {
			s.Skip()
		}
		// - 場札がドロー2の時、ドロー2の効果が適用される。1番目のプレーヤーが2枚引き、2番目のプレイヤーからスタート。
		if card.IsDraw() {
			s.ApplyDrawEffect()
		}
		// - 場札がリバースの時、リバースの効果が適用される。3番目のプレーヤーからリバースしてスタート。
		if card.IsDraw() {
			s.Reverse()
			s.Next()
		}
		// - 場札が決まった後に手札を配る。(公式ルールとは山札の決定と手札の配る順番が異なる。)
		for i := 0; i < len(s.Players); i++ {
			ok := DealCards(i, 7, s)
			if !ok {
				// 山札も捨て札もない。全てのカードが手札にある???
				// このタイミングで起きるのはおかしい。
				panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???"))
			}
		}
	} else {
		// 公式ルール。
		// 場からプレイヤーの手札を配る。初期7枚。
		for i := 0; i < len(s.Players); i++ {
			ok := DealCards(i, 7, s)
			if !ok {
				// 山札も捨て札もない。全てのカードが手札にある???
				// このタイミングで起きるのはおかしい。
				panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???"))
			}
		}

		if Developing && Debug {
			log.Printf("Deck   : %#v", PrintCards(s.Deck))
			log.Printf("Deck Top5: %#v", PrintTopCards(s.Deck, 5))
			log.Printf("Discard: %#v", PrintCards(s.Discard))
		}

		// 場札の決定
		{
			/* 説明書より
			   1. 最初に親を決めて、カードを切ります。
			   2. 親は、各プレイヤーに7枚ずつカードを伏せて配ります。
			   3. 残りのカードは、伏せて積んでおきます。これが引き札の山です。
			   4. 親が引き札の山の一番上のカードを1枚めくり、脇に置きます。
			      これが捨て札の山の最初の1枚になります。捨て札の山の一番上のカードを場のカードと呼びます。
			      最初の場のカードが記号カードだった時は、引き札の山の中に戻して、次の1枚を場のカードにします。
			   5. 親の左どなりの人が最初のプレイヤーです。時計周りにカードを捨てていきます。
			*/

			// 4. 親が引き札の山の一番上のカードを1枚めくり、脇に置きます。
			//    これが捨て札の山の最初の1枚になります。捨て札の山の一番上のカードを場のカードと呼びます。
			//    最初の場のカードが記号カードだった時は、引き札の山の中に戻して、次の1枚を場のカードにします。

			for {
				c, ok := pop(&s.Deck)
				if !ok {
					panic(errors.Errorf("%v", "バグ"))
				}
				s.LastCard = &c
				if s.LastCard.IsAction() {
					InsertRandomCard(&s.Deck, *s.LastCard)
				} else {
					push(&s.Discard, *s.LastCard)
					break
				}
			}
		}
	}

	return s
}

func (s *State) Print() string {
	num := ""
	if s.Number >= 0 {
		num = fmt.Sprintf("%v", s.Number)
	}
	turn := "CW"
	if s.Turn < 0 {
		turn = "CCW"
	}
	return fmt.Sprintf("(%v)%v%v+%v%v",
		s.Players[s.NowPlayerIndex], // 今の人のID
		ColorName(s.Color),          // -1の時任意
		num,                         // -1の時任意
		s.Draw,                      // 0以下の時任意
		turn,                        // CW:1 CCW:-1
	)
}
func (s *State) PrintDetail() string {
	return fmt.Sprintf(`
	#############################################################
	Color          : %10v : -1の時任意
	Number         : %10v : -1の時任意
	Draw           : %10v : 0以下の時任意
	Turn           : %10v : CW:1 CCW:-1
	NowPlayerIndex : %10v : 今の人のID
	Players        : %v : 参加者
	NPCs           : %v : NPCかどうか
	Hands     len  : %10v : 手札
	Deck      len  : %10v : 山札
	Discard   len  : %10v : 捨て札
	SkipCount      : %10v : 全員手札が出せない状況になったら引き分けにする。
	#############################################################
	`,
		ColorName(s.Color),
		s.Number,
		s.Draw,
		s.Turn,
		s.NowPlayerIndex,
		s.Players,
		s.NPCs,
		len(s.Hands),
		len(s.Deck),
		len(s.Discard),
		s.SkipCount,
	)
}

func (s *State) Reverse() {
	fmt.Printf("\t> 順番が反転します。\n")
	if s.Turn > 0 {
		s.Turn = -1
	} else {
		s.Turn = 1
	}
}
func (s *State) Skip() {
	fmt.Printf("\t> %v さんはスキップされます。\n", s.NextName())
	s.NowPlayerIndex = ((s.NowPlayerIndex + s.Turn*2) + len(s.Players)*2) % len(s.Players)
}
func (state *State) DiscardAll(c *Card) {
	// R.official.08 ディスカードオールカードの追加枚数。各色枚数。
	fmt.Printf("\t> 手札の %v のカード全てを捨てます。\n", c.ColorName())

	// 捨てられたディスカードオールを最終的に捨札の一番上になるように保持しておく。
	card := *c
	// 捨てられたディスカードオールを手札から取り除く。
	*state.Hand() = RemoveCard(*state.Hand(), card)
	// 捨てられたディスカードオールと同じ色のカードを捨札に捨てる。
	tmp := make([]Card, 0, len(*state.Hand())) // 最終的な手札を作る。
	for _, h := range *state.Hand() {
		if card.Color() == h.Color() {
			push(&state.Discard, h) // 捨札に送る。
		} else {
			tmp = append(tmp, h)
		}
	}
	push(&state.Discard, card) // ディスカードオールを最後に捨てる。
	*state.Hand() = tmp        // 手札を設定
}
func (s *State) ApplyDrawEffect() {
	if s.Draw > 0 {
		fmt.Printf("\t> %v さんは%v枚引きます。\n", s.NowPlayerIndex, s.Draw)
		ok := DealCards(s.NowPlayerIndex, s.Draw, s)
		if !ok {
			s.SkipCount++
			if Debug {
				fmt.Printf("\t> 山札がないので %v さんはカードを引くことができませんでした。 -- SkipCount: %v\n", s.Name(), s.SkipCount)
			} else {
				fmt.Printf("\t> 山札がないので %v さんはカードを引くことができませんでした。\n", s.Name())
			}
			if s.SkipCount > len(s.Players)*2 {
				// 誰も出せない
				fmt.Printf("> 全員出せるカードが無かったため引き分けです。\n")
				os.Exit(-1)
			}
		} else {
			s.SkipCount = 0
		}
		s.Draw = 0
		s.Skip()
	}
}
func (s *State) Next() {
	s.NowPlayerIndex = ((s.NowPlayerIndex + s.Turn) + len(s.Players)*2) % len(s.Players)
}
func (s *State) Update() {
	s.Next()
	s.ApplyDrawEffect()
}
func (s *State) Name() string {
	return s.Players[s.NowPlayerIndex]
}
func (s *State) NextPlayerID() int {
	return ((s.NowPlayerIndex + s.Turn) + len(s.Players)) % len(s.Players)
}
func (s *State) NextName() string {
	return s.Players[s.NextPlayerID()]
}
func (s *State) OtherPlayerNames() []string {
	o := []string{}
	for _, n := range s.Players {
		if n != s.Name() {
			o = append(o, n)
		}
	}
	return o
}
func (s *State) OtherHandCounts() []int {
	o := []int{}

	for i := 0; i < len(s.Players); i++ {
		if i != s.NowPlayerIndex {
			o = append(o, len(s.Hands[i]))
		}
	}
	return o
}
func (s *State) Hand() *[]Card {
	return &s.Hands[s.NowPlayerIndex]
}
func (s *State) HandID() []int {
	hs := s.Hand()
	ret := []int{}
	for _, h := range *hs {
		ret = append(ret, h.Type)
	}
	return ret
}
