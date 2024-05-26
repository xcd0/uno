package core

import (
	"bufio"
	crypto_rand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func ChooseCard(state *State, npc bool) *Card {
	var chosen Card
	tmp := FilterCards(state)
	// forの間一時的にコピーしたものから取り出す。
	dealed := false
	for {
		if len(tmp) == 0 {
			if !dealed {
				dealed = true
				// 出せるカードがないので追加で1枚引く。
				fmt.Printf("\t> %v さんは出せるカードがないので1枚引きます。\n", state.Name())
				ok := DealCards(state.NowPlayerIndex, 1, state)
				if !ok {
					fmt.Printf("\t> 山札がないので %v さんはカードを引くことができませんでした。スキップします。\n", state.Name())
					return nil
				}
				continue
			} else {
				return nil
			}
		}
		if Developing && Debug {
			fmt.Printf("\t> %v さんの手札で出せるもの : %v\n", state.Name(), PrintCardsDescription(tmp))
		}

		if selected := SelectCard(state, tmp, npc); selected == nil {
			return nil
		} else {
			chosen = *selected
		}

		if CheckCard(chosen, state) {
			if !npc || Developing && Debug {
				fmt.Printf(" > OK!\n")
			}
			break
		}
		if !npc || Developing && Debug {
			fmt.Printf(" > NG!\n")
		}

	}
	*state.Hand() = RemoveCard(*state.Hand(), chosen)
	return &chosen
}

func FilterCards(state *State) []Card {
	ret := make([]Card, 0, len(*state.Hand()))
	for _, c := range *state.Hand() {
		if CheckCard(c, state) {
			if Developing && Debug {
				log.Printf("%v -> %#v OK", state.Print(), c.Print())
			}
			ret = append(ret, c)
		} else {
			if Developing && Debug {
				log.Printf("%v -> %#v NG", state.Print(), c.Print())
			}
		}
	}
	if Developing && Debug {
		log.Printf("ret: %v", PrintCards(ret))
		//log.Printf("%v", ret)
	}
	return ret
}

func SelectCard(state *State, tmp []Card, npc bool) *Card {
	var chosen Card
	var ok bool
	if npc {
		// ランダムに1枚選ぶ実装
		chosen, ok = PopRandomCard(&tmp)
		if !ok {
			panic(errors.Errorf("%v", "バグ"))
		}
	} else if true {
		reader := bufio.NewReader(os.Stdin)
		dealed := false // 1回だけ引ける。引いたら引いたカードしか出せない。
		for {
			min_num := 0
			fmt.Printf("\t\t> 捨てるカードを選択してください。\n")
			// 出せるカードのリストを提示。
			if len(state.Deck) == 0 {
				min_num = 1
			}
			if !dealed {
				fmt.Printf("\t\t\t> 0: 出さずに山札から1枚引く。\n") // 1回だけ引ける。引いたら引いたカードしか出せない。
			} else {
				fmt.Printf("\t\t\t> 0: スキップする。\n")
				fmt.Printf("\t\t\t> 1: 引いたカードを出す。\n")
			}
			// ハウスルール
			// state.Rule.AllowSequentialSameColor         // 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
			// state.Rule.AllowSimultaneousAnyColor        // 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
			// state.Rule.AllowWinWithoutUnoOnSimultaneous // 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可

			if !npc {
				for i, c := range tmp {
					fmt.Printf("\t\t\t> %v: %v (%v) %v\n", i+1, c.Description, c.Name, c.Detail)
				}
			}

			max_num := len(tmp)
			var input string
			var err error
			if npc || Developing && Debug {
				// デバッグ中はとりあえずランダムに選択。
				num := 0
				if min_num == 0 {
					num = rand.Intn(max_num+1) - 1 // 0 - max_num
				} else {
					num = rand.Intn(max_num) + 1 // 1 - max_num
				}
				//num = min(max(num, max_num), min_num)
				input = fmt.Sprintf("%v", num) // 0-max_numのはず。
				fmt.Printf("%v ", input)
			} else {
				fmt.Printf("\t\t\t\t> ")
				input, err = ReadRune(reader)
				fmt.Printf("%v ", input)
				//log.Printf("input %v", input)
			}
			if err != nil {
				if Developing && Debug {
					log.Printf("\n\t\t> %vから%vの数を指定して下さい。", min_num, max_num)
				}
				fmt.Printf("\n\t\t> %vから%vの数を指定して下さい。\n", min_num, max_num)
				continue
			}
			input = strings.TrimSpace(input)   // 改行文字を取り除く
			number, err := strconv.Atoi(input) // 入力を整数に変換
			if err != nil {
				if Developing && Debug {
					log.Printf("\n\t\t> %vから%vの数を指定して下さい。", min_num, max_num)
				}
				fmt.Printf("\n\t\t> %vから%vの数を指定して下さい。\n", min_num, max_num)
				continue
			}
			if number == 0 {
				// 山札から引く。
				if len(state.Deck) == 0 {
					fmt.Printf("\n\t\t> 山札がありません。\n", min_num, max_num)
					fmt.Printf("\n\t\t> 1から%vの数を指定して下さい。\n", min_num, max_num)
				} else {
					dealed = true
					ok := DealCards(state.NowPlayerIndex, 1, state)
					if !ok {
						log.Printf("state : %v", state.Print())
						log.Printf("%v", state.PrintDetail())
						panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???"))
					}
					c := peek(&state.Hands[state.NowPlayerIndex])
					fmt.Printf("\n\t\t> %v (%v) を引きました。\n", c.Description, c.Name)
					if state.Rule.DisallowImmediatePlayOfDrawCards { // 引いた直後のドロー2またはドロー4をすぐに出すことを禁止する。
						if action := c.ActionType() + 100; action == DRAW2__ && action == DRAW2__ {
							fmt.Printf("\t\t\t> 引いた直後のドロー2またはドロー4をすぐに出すことを禁止するルールが適用されています。\n")
							fmt.Printf("\t\t\t> 引いたカードが出せないためスキップします。\n")
							return nil
						}
					}
					if !CheckCard(peek(&state.Hands[state.NowPlayerIndex]), state) {
						fmt.Printf("\t\t\t> 引いたカードが出せないためスキップします。\n")
						return nil
					}
					// 引いた場合、引いたカードのみだせる。元の手札空は出せない。

					if state.Rule.PlayFromHandAfterDraw { // (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。
						tmp = append(tmp, peek(&state.Hands[state.NowPlayerIndex]))
					} else {
						tmp = []Card{peek(&state.Hands[state.NowPlayerIndex])}
					}
					continue
				}
			} else if number > min_num && number <= max_num { // 数値が1から4の間かどうかをチェック
				chosen = tmp[number-1]
				fmt.Printf(" > %v を選択しました。", chosen.Print())
				break // ok
			}

			if Developing && Debug {
				log.Printf("\n\t\t> %vから%vの数を指定して下さい。", min_num, max_num)
			}
			fmt.Printf("\n\t\t> %vから%vの数を指定して下さい。\n", min_num, max_num)
		}
	}
	return &chosen
}

func CheckCard(chosen Card, state *State) bool {
	// 選んだカードのチェック
	if chosen.IsNum() {
		if chosen.Color() == state.Color {
			if Developing && Debug {
				log.Printf("c:%v %v, s.c:%v -> %v", chosen.Print(), chosen.Color(), state.Color, chosen.Color() == state.Color)
			}
			return true
		}
		if chosen.Num() == state.Number {
			if Developing && Debug {
				log.Printf("c:%v %v, s.c:%v -> %v", chosen.Print(), chosen.Color(), state.Color, chosen.Num() == state.Number)
			}
			return true
		}
		if Developing && Debug {
			log.Printf("c:%v %v, s.c:%v -> false", chosen.Print(), chosen.Color(), state.Color)
		}
	}
	if chosen.IsAction() {
		action := chosen.ActionType() + 100
		if Developing && Debug {
			log.Printf("chosen.Action: %v, action: %v", chosen.ActionType(), action)
		}
		switch action {
		case WILD___: // 場のカードがどのカードであっても出せる。
			if state.LastCard == nil { // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
				return true
			}
			if Developing && Debug {
				l := peek(&state.Discard)
				log.Printf("c:%v --, s:%v -> true", chosen.Print(), l.Print())
			}
			return true
		case DRAW4__:
			if state.LastCard == nil { // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
				return true //
			} else {
				// スキップされていないとき。
				// チャレンジルールはここではないところで判定する。
				// 前の人のカードがドローカードの時、ルール毎に出せるかどうか決まる。
				discard := peek(&state.Discard)
				action := discard.ActionType() + 100
				if action == DRAW2__ && action == DRAW2__ {
					if state.Rule.RequireMatchingDrawTwoForChallenge && action == DRAW2__ { // R.jua.05      ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
						if Developing && Debug {
							log.Printf("RequireMatchingDrawTwoForChallenge: %v : R.jua.05 ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。", state.Rule.RequireMatchingDrawTwoForChallenge)
						}
						return state.Color == discard.Color()
					}
					return state.Rule.StackDrawCards // R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
				}
				return true // ドローカード以外は場のカードがどのカードであっても出せる。
			}
		case DRAW2__:
			fallthrough
		case REVERSE:
			fallthrough
		case SKIP___:
			fallthrough
		case DISCARD:
			if state.Color == -1 {
				return true
			}
			return chosen.Color() == state.Color // 色が一致する必要がある。
		default:
			// 来ないはず。
			if Developing && Debug {
				log.Printf("chosen: %v", chosen)
			}
			panic(errors.Errorf("%v", "バグ"))
		}
	}
	return false // だめ。
}

func PopRandomCard(cards *[]Card) (Card, bool) {
	n := len(*cards)
	if n == 0 {
		return Card{}, false // スライスが空の場合は空のCardとnilを返す
	}

	// ランダムなインデックスを選択
	indexBig, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic("Failed to generate random index: " + err.Error())
	}
	index := int(indexBig.Int64())

	// 選ばれたカードを取り出す
	selectedCard := (*cards)[index]

	// 選ばれたカードをスライスから削除
	// スライスの末尾の要素と選ばれたインデックスの要素を交換し、スライスをリサイズ
	(*cards)[index] = (*cards)[n-1]
	(*cards) = (*cards)[:n-1]

	return selectedCard, true

}

// ReadRune は標準入力から1文字読み込み、その文字とエラーを返します。
func ReadRune(reader *bufio.Reader) (string, error) {
	// ターミナルの設定を変更して入力をすぐに受け取る
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer func() {
		// プログラム終了時に元のターミナル設定を復元
		exec.Command("stty", "-F", "/dev/tty", "sane").Run()
	}()
	char, _, err := reader.ReadRune()
	if err != nil {
		return "", err // エラー発生時には0とエラー情報を返す
	}
	return string(char), nil // 正常に読み込めた場合、文字とnilを返す
}
