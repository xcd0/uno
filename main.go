package main

import (
	"fmt"
	"log"
	math_rand "math/rand"
	"time"

	"github.com/pkg/errors"
)

var (
	version         string = "debug"
	revision        string
	ThisProgramPath string // $0の絶対パスを入れる
	CurrentPath     string // 起動時のカレントディレクトリ

	debug      = true
	developing = false
)

func main() {
	args, s := ArgParse()
	Run(args, s)
}

func Run(args *Args, s *Setting) {
	math_rand.Seed(time.Now().UnixNano())

	// 使用するルール選択
	// とりあえず公式
	rule := (&UnoRule{}).Init()
	rule.SetRule(UnoRuleList[0])

	state := InitGame(
		(State{}).Init(
			args.PlayerNames,
			args.NumberOfPlayers,
			GetCards(rule),
			&rule,
		),
	)

	setting := ReadSetting(args, NewSetting(args))

	switch {
	case args.ArgsServer != nil: // サーバーモードとして起動する。
		go UnoServer(setting.Port, rule)
	case args.ArgsClient != nil: // クライアントモードとして起動する。
		UnoClient(setting.Port)
	case args.ArgsServerClient != nil: // サーバーとクライアント同時に起動する。
		go UnoServer(setting.Port, rule)
		UnoClient(setting.Port)
	case args.ArgsSolo != nil: // 一人プレイモードで起動する。
		go UnoServer(setting.Port, rule)
		UnoClient(setting.Port)
	case args.CreateEmptyHjson != nil: // 空の設定ファイルを生成する。
	case args.ConvertToJson != nil: // 設定ファイルをjsonに変換する。
	case args.VersionSub != nil: // バージョン番号を出力する。-vと同じ。
		//case args.Readme: // readme.mdを出力
		//fmt.Printf("%s", "説明文を書く?")
	}

	return
	GameStart(state)
	return
}

func InitGame(state *State) *State {

	// shuffle
	ShuffleCards(state.Deck)
	CryptoRandShuffle(state.Deck)

	// 場からプレイヤーの手札を配る。初期7枚。
	for i := 0; i < len(state.Players); i++ {
		ok := DealCards(i, 7, state)
		if !ok {
			// 山札も捨て札もない。全てのカードが手札にある???
			// このタイミングで起きるのはおかしい。
			panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???"))
		}
	}

	if developing && debug {
		log.Printf("Deck   : %#v", PrintCards(state.Deck))
		log.Printf("Deck Top5: %#v", PrintTopCards(state.Deck, 5))
		log.Printf("Discard: %#v", PrintCards(state.Discard))
	}

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
		c, ok := pop(&state.Deck)
		if !ok {
			panic(errors.Errorf("%v", "バグ"))
		}
		state.LastCard = &c
		if state.LastCard.IsAction() {
			InsertRandomCard(&state.Deck, *state.LastCard)
		} else {
			push(&state.Discard, *state.LastCard)
			break
		}
	}

	return state
}

func GameStart(state *State) {

	// 5. 親の左隣の人が最初のプレイヤーです。時計周りにカードを捨てていきます。
	fmt.Printf("> %v さんが親です。親の次の人が最初です。\n", state.Name())
	fmt.Printf("> 最初のカードは %v です。\n", state.Name)
	state.Update(state)
	fmt.Printf("> %v さんのターンです。\n", state.Name())

	//for i := 0; i < 3; i++ {
	for {
		// 現在のプレーヤーがNPCか判断
		npc := state.NowPlayerIndex != 0 // とりあえず実装

		if !npc || developing && debug {
			fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), PrintCardsDescription(*state.Hand()))
		}

		// 捨てるカードの選択
		chosen := ChooseCard(state, npc)
		if chosen == nil {
			fmt.Printf("\t> %v さんは出せるカードがありません。スキップします。\n", state.Name())
			state.SkipCount++
			if state.SkipCount > state.Rule.EveryoneLosesAfterConsecutiveSkips*len(state.Players) { // R.algori.O5 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
				fmt.Printf("\t> %v 週スキップしたため全員敗者とします。\n", state.Rule.EveryoneLosesAfterConsecutiveSkips)
				break
			}
			TurnEnd(state, npc)
			continue
		}

		// 捨てたカードがアクションカードならアクションカードの効果だけ発動させる。

		if chosen.IsNum() {
			//if developing && debug {
			//	log.Println("数値カード")
			//}
			state.Color = chosen.Color()
			state.Number = chosen.Num()
		} else if chosen.IsAction() {
			//if developing && debug {
			//	log.Println("アクションカード")
			//}
			if len(*state.Hand()) != 1 { // 残り1枚の時勝ち確定なので不要
				ActionEffect(chosen, state, npc)
			}
		}

		{
			if len(*state.Hand()) == 1 {
				fmt.Printf("\t> %v さんは %v を捨てました。UNO! です。\n", state.Name(), chosen.Print())
			} else {
				fmt.Printf("\t> %v さんは %v を捨てました。\n", state.Name(), chosen.Print())
			}

			if len(*state.Hand()) == 0 {
				if state.Rule.RestrictWinOnSpecialCard { // R.house-jp.02 記号カードでの勝利を禁止する。
					fmt.Printf("\t\t> アクションカードでの勝利を禁止するルールが適用されています。ペナルティとして2枚ドローします。\n")
					// 手札に戻す。
					push(state.Hand(), *chosen)
					if !DealCards(state.NowPlayerIndex, 2, state) {
						fmt.Printf("\t\t> 山札がないためドローできませんでした。\n")
					}
					TurnEnd(state, npc)
				}
				fmt.Printf("\t> %v さんの手札が無くなりました。%v さんの勝ちです。\n", state.Name(), state.Name())
				break
			} else if len(*state.Hand()) == 1 {
				//fmt.Printf("\t> %v さんは UNO! です。\n", state.Name())
			} else {
				if developing && debug {
					fmt.Printf("\t> %v さんの残り枚数は %v です。\n", state.Name(), len(*state.Hand()))
				}
			}
		}

		TurnEnd(state, npc)
	}
}

func TurnEnd(state *State, npc bool) {
	if !npc || developing && debug {
		fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), PrintCardsDescription(*state.Hand()))
	}

	hands := ""
	for i, h := range state.Hands {
		if i != 0 {
			hands += ","
		}
		//log.Printf("%v %v", i, state.NowPlayerIndex)
		if i == state.NowPlayerIndex {
			hands += fmt.Sprintf("<%v>", len(h))
		} else {
			hands += fmt.Sprintf("%v", len(h))
		}
	}
	fmt.Printf("\t---------------------------------> %v (%v)\n", state.Print(), hands)
	state.Update(state)
	fmt.Printf("> %v さんのターンです。\n", state.Name())
}
