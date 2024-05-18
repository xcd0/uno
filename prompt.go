package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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
			if npc || developing && debug {
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
				if developing && debug {
					log.Printf("\n\t\t> %vから%vの数を指定して下さい。", min_num, max_num)
				}
				fmt.Printf("\n\t\t> %vから%vの数を指定して下さい。\n", min_num, max_num)
				continue
			}
			input = strings.TrimSpace(input)   // 改行文字を取り除く
			number, err := strconv.Atoi(input) // 入力を整数に変換
			if err != nil {
				if developing && debug {
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

			if developing && debug {
				log.Printf("\n\t\t> %vから%vの数を指定して下さい。", min_num, max_num)
			}
			fmt.Printf("\n\t\t> %vから%vの数を指定して下さい。\n", min_num, max_num)
		}
	}
	return &chosen
}

func SelectColor(state *State, npc bool) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\t\t> 色を指定してください。:\n")
		fmt.Printf("\t\t> 青:1 / 緑:2 / 赤:3 / 黄:4 > ")

		var input string
		var err error

		if npc || developing && debug {
			// デバッグ中はとりあえずランダムに選択。
			input = fmt.Sprintf("%v", rand.Intn(4)+1) // 1-4のはず。
			fmt.Printf("%v ", input)
		} else {
			input, err = ReadRune(reader)
		}

		if err != nil {
			if developing && debug {
				log.Printf("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
			continue
		}
		input = strings.TrimSpace(input)   // 改行文字を取り除く
		number, err := strconv.Atoi(input) // 入力を整数に変換
		if err != nil {
			if developing && debug {
				log.Printf("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
			continue
		}
		// 数値が1から4の間かどうかをチェック
		if number >= 1 && number <= 4 {
			switch number {
			case 1:
				state.Color = BLUE__
			case 2:
				state.Color = GREEN_
			case 3:
				state.Color = RED___
			case 4:
				state.Color = YELLOW
			}
			color := ColorName(state.Color)
			fmt.Printf("> 次の色は %v です。\n", color)
			break
		}
		if developing && debug {
			log.Printf("\t\t> 1から4の数を指定して下さい。")
		}
		fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
	}
}
