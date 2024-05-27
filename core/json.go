package core

// 参加者がアクセスできるゲームの情報。
type JsonGameState struct {
	CurrentCard      int      `json:"current_card"             comment:"場に出ているカードのID。場に最後に捨てられたカードのID。"`
	DrawnCardID      []int    `json:"drawn_card_id"            comment:"新規に引いたカードのID。draw2等で複数引いている場合がある。"`
	Hand             []int    `json:"player_hand"              comment:"自分の手札。新規に引いたカードも含む。"`
	TeammateHand     []int    `json:"teammate_player_hand"     comment:"仲間の手札。4人を2人づつのペアにしてペアで勝利を競うルールがある。対角のプレーヤーをチームメイトとするルールの場合、対角のプレーヤーの手札を見られる。"`
	YourName         string   `json:"your_name"                comment:"自分の名前。"`
	OtherPlayerNames []string `json:"other_player_names"       comment:"他プレイヤーの名前のリスト。このリストの要素番号が、クライアントからサーバーへ返すときのプレイヤーのIDになる。"`
	OtherHandCounts  []int    `json:"other_player_hand_counts" comment:"他プレイヤーの手札枚数のリスト。"`
	UnoCalled        []bool   `json:"uno_called"               comment:"他プレイヤーがUNO!とコールしたかどうかリスト。"`
	Direction        int      `json:"direction"                comment:"順番の回る方向。1の時他プレーヤーのリストの先頭、-1の時末尾のプレーヤーが次になる。"`
	RuleToApply      string   `json:"rule"                     comment:"使用されているルール。"default:"official_new"`
	CustomRules      UnoRule  `json:"custom-rules"             comment:"ユーザーがルールをカスタマイズできる。"`
	Port             int      `json:"port"                     comment:"サーバーが使用するポート番号" default:"5000"`
}

func (jgs *JsonGameState) Init(yourname string, state *State) JsonGameState {
	others := make([]string, 0, len(state.Players)-1)
	uno_called := make([]bool, 0, len(state.Players)-1)
	other_hands := make([]int, 0, len(state.Players)-1)
	for _, n := range state.Players {
		if n != yourname {
			others = append(others, n)
			uno_called = append(uno_called, false)
			other_hands = append(other_hands, 0)
		}
	}

	*jgs = JsonGameState{
		CurrentCard:      -1,                        // 場に出ているカードのID。場に最後に捨てられたカードのID。
		DrawnCardID:      make([]int, 0, 10),        // 新規に引いたカードのID。draw2等で複数引いている場合がある。
		Hand:             make([]int, 0, 30),        // 自分の手札。新規に引いたカードも含む。
		TeammateHand:     make([]int, 0, 30),        // 仲間の手札。4人を2人づつのペアにしてペアで勝利を競うルールがある。対角のプレーヤーをチームメイトとするルールの場合、対角のプレーヤーの手札を見られる。
		YourName:         yourname,                  // 自分の名前。
		OtherPlayerNames: others,                    // 他プレイヤーの名前のリスト。このリストの要素番号が、クライアントからサーバーへ返すときのプレイヤーのIDになる。
		OtherHandCounts:  other_hands,               // 他プレイヤーの手札枚数のリスト。
		UnoCalled:        uno_called,                // 他プレイヤーがUNO!とコールしたかどうかリスト。
		Direction:        1,                         // 順番の回る方向。1の時他プレーヤーのリストの先頭、-1の時末尾のプレーヤーが次になる。
		RuleToApply:      state.Setting.RuleToApply, // 使用されているルール。
		CustomRules:      state.Rule,                //ユーザーがルールをカスタマイズできる。
		Port:             state.Setting.Port,        //サーバーが使用するポート番号。
	}
	return *jgs
}

type JsonClientPlay struct {
	Discard           []int `json:"discard_id"`            // 捨てるカードを指定する。空の時、初回に限り山札からカードを1枚引いて再度カードを選びなおす事ができる。
	CallUno           bool  `json:"call_uno"`              // 基本falseを返す。自分がUNO!とコールするときtrueにする。
	CallUnoOnPlayerID int   `json:"call_uno_on_player_id"` // 基本-1を返す。他人がUNO!とコールしていない時指摘する場合にそのプレーヤーIDを指定する。
}

func (jcp *JsonClientPlay) Init() JsonClientPlay {
	*jcp = JsonClientPlay{
		Discard:           make([]int, 0, 10), // []int `json:"discard_id"`            // 捨てるカードを指定する。空の時、初回に限り山札からカードを1枚引いて再度カードを選びなおす事ができる。
		CallUno:           false,              // bool  `json:"call_uno"`              // 基本falseを返す。自分がUNO!とコールするときtrueにする。
		CallUnoOnPlayerID: -1,                 // int   `json:"call_uno_on_player_id"` // 基本-1を返す。他人がUNO!とコールしていない時指摘する場合にそのプレーヤーIDを指定する。
	}
	return *jcp
}

// 全ての使用するカードの情報をクライアントがサーバーに問い合わせることができる。
type JsonCardInfo struct {
	CardAll []Card `json:"card_info"`
}

func (jci *JsonCardInfo) Init(rule UnoRule) JsonCardInfo {
	jci.CardAll = GetCards(rule)
	return *jci
}
