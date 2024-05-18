package main

import (
	"fmt"
	"math/rand"
	"os"
)

type State struct {
	Color          int      // -1の時任意
	Number         int      // -1の時任意
	Draw           int      // 0以下の時任意
	Turn           int      // CW:1 CCW:-1
	NowPlayerIndex int      // 今の人のID
	Players        []string // 参加者
	NPCs           []bool   // NPCかどうか
	Hands          [][]Card // 手札
	Deck           []Card   // 山札
	Discard        []Card   // 捨て札
	LastCard       *Card    // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
	SkipCount      int      // 全員手札が出せない謎状況になったら引き分けにする。
	Rule           UnoRule  // 使用中のルール
}

func (state State) Init(players []string, numberOfPlayers int, cards []Card, rule *UnoRule) *State {

	if len(players) == 0 {
		players = []string{"You", "P1", "P2", "P3"}
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
	}
	return &state
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
func (s *State) Update(state *State) {
	s.NowPlayerIndex = ((s.NowPlayerIndex + s.Turn) + len(s.Players)) % len(s.Players)
	if s.Draw > 0 {
		fmt.Printf("\t> %v さんは%v枚引きます。\n", s.NowPlayerIndex, s.Draw)
		ok := DealCards(s.NowPlayerIndex, s.Draw, state)
		if !ok {
			s.SkipCount++
			if debug {
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
func (s *State) Name() string {
	return s.Players[s.NowPlayerIndex]
}
func (s *State) NextName() string {
	tmp := ((s.NowPlayerIndex + s.Turn) + len(s.Players)) % len(s.Players)
	return s.Players[tmp]
}
func (s *State) Hand() *[]Card {
	return &s.Hands[s.NowPlayerIndex]
}
