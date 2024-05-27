package main

import (
	"fmt"
	"io"
	math_rand "math/rand"
	"os"
	"time"

	"github.com/pkg/errors"
	_ "github.com/rivo/tview"
	"github.com/xcd0/uno/client"
	"github.com/xcd0/uno/core"
	"github.com/xcd0/uno/server"
)

var (
	version         string = "Debug"
	revision        string
	ThisProgramPath string // $0の絶対パスを入れる
	CurrentPath     string // 起動時のカレントディレクトリ
	wrapperStdout   io.Writer
	wrapperStderr   io.Writer
	logfile         *os.File
)

func Init() {
	math_rand.Seed(time.Now().UnixNano())
}

func main() {
	args := ArgParse()
	Run(args)
}

func Run(args *Args) {
	// 使用するルール選択
	// TODO: 実装
	// とりあえず公式
	rule := (&core.UnoRule{}).Init()
	rule.SetRule(core.UnoRuleList[0].Name)

	state := core.NewState(
		args.PlayerNames,
		args.NumberOfPlayers,
		core.GetCards(rule),
		&rule,
	)
	state.ShuffleDeck()
	state.DealHandsOnInit()

	if !true {
		// TUIクライアントとりあえず実装。
		client.Tui()
	} else if !true {
		// とりあえず実装。
		GameStart(state)
		return
	} else {

		// サーバーとクライアントを分けた実装。

		{
			setting := core.ReadSetting(args.SettingPath, core.NewSetting())
			loggingSettings(args.LogPath, &setting)
			if args.CreateEmptyHjson != nil {
				core.CreateEmptyHjson()
				os.Exit(0)
			}
			state.Setting = &setting
		}

		switch {
		case args.ArgsServer != nil: // サーバーモードとして起動する。
			server.UnoServer(state)
		case args.ArgsClient != nil: // クライアントモードとして起動する。
			client.UnoClient(state.Setting.Port)
		case args.ArgsServerClient != nil: // サーバーとクライアント同時に起動する。
			go server.UnoServer(state)
			client.UnoClient(state.Setting.Port)
		case args.ArgsSolo != nil: // 一人プレイモードで起動する。
			go server.UnoServer(state)
			client.UnoClient(state.Setting.Port)
		case args.CreateEmptyHjson != nil: // 空の設定ファイルを生成する。
		// TODO: 実装
		case args.ConvertToJson != nil: // 設定ファイルをjsonに変換する。
		// TODO: 実装
		case args.VersionSub != nil: // バージョン番号を出力する。-vと同じ。
			// TODO: 実装
			//case args.Readme: // readme.mdを出力
			//fmt.Printf("%s", "説明文を書く?")
			// TODO: 実装
		}
		return
	}
}

func GameStart(state *core.State) {

	// 5. 親の左隣の人が最初のプレイヤーです。時計周りにカードを捨てていきます。
	fmt.Printf("> %v さんが親です。親の次の人が最初です。\n", state.Name())
	fmt.Printf("> 最初のカードは %v です。\n", state.Name)
	state.Update()
	fmt.Printf("> %v さんのターンです。\n", state.Name())

	//for i := 0; i < 3; i++ {
	for {
		// 現在のプレーヤーがNPCか判断
		npc := state.NowPlayerIndex != 0 // とりあえず実装

		if !npc || core.Developing && core.Debug {
			fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), core.PrintCardsDescription(*state.Hand()))
		}

		// 捨てるカードの選択
		chosen := core.ChooseCard(state, npc)
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
			//if core.Developing && core.Debug {
			//	log.Println("数値カード")
			//}
			state.Color = chosen.Color()
			state.Number = chosen.Num()
		} else if chosen.IsAction() {
			//if core.Developing && core.Debug {
			//	log.Println("アクションカード")
			//}
			if len(*state.Hand()) != 1 { // 残り1枚の時勝ち確定なので不要
				core.ActionEffect(chosen, state, npc)
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
					if !core.DealCards(state.NowPlayerIndex, 2, state) {
						fmt.Printf("\t\t> 山札がないためドローできませんでした。\n")
					}
					TurnEnd(state, npc)
				}
				fmt.Printf("\t> %v さんの手札が無くなりました。%v さんの勝ちです。\n", state.Name(), state.Name())
				break
			} else if len(*state.Hand()) == 1 {
				//fmt.Printf("\t> %v さんは UNO! です。\n", state.Name())
			} else {
				if core.Developing && core.Debug {
					fmt.Printf("\t> %v さんの残り枚数は %v です。\n", state.Name(), len(*state.Hand()))
				}
			}
		}

		TurnEnd(state, npc)
	}
}

func TurnEnd(state *core.State, npc bool) {
	if !npc || core.Developing && core.Debug {
		fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), core.PrintCardsDescription(*state.Hand()))
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
	state.Update()
	fmt.Printf("> %v さんのターンです。\n", state.Name())
}

// push関数：任意の型の要素を受け取るためにCardを使用
func push(slice *[]core.Card, value core.Card) {
	*slice = append(*slice, value)
}

// pop関数：スライスから要素を削除
func pop(slice *[]core.Card) (core.Card, bool) {
	if len(*slice) == 0 {
		return core.Card{}, false // スライスが空の場合は何もしない
	}
	// スライスの最後の要素を取得し、スライスからその要素を除去
	value := (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
	return value, true
}

// peek関数：スタックの一番上を覗く
func peek(slice *[]core.Card) core.Card {
	if len(*slice) == 0 {
		panic(errors.Errorf("%v", "バグ"))
	}
	return (*slice)[len(*slice)-1]
}
