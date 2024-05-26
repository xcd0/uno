package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/xcd0/uno/core"
)

func SelectColor(state *core.State, npc bool) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\t\t> 色を指定してください。:\n")
		fmt.Printf("\t\t> 青:1 / 緑:2 / 赤:3 / 黄:4 > ")

		var input string
		var err error

		if npc || core.Developing && core.Debug {
			// デバッグ中はとりあえずランダムに選択。
			input = fmt.Sprintf("%v", rand.Intn(4)+1) // 1-4のはず。
			fmt.Printf("%v ", input)
		} else {
			input, err = core.ReadRune(reader)
		}

		if err != nil {
			if core.Developing && core.Debug {
				log.Printf("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
			continue
		}
		input = strings.TrimSpace(input)   // 改行文字を取り除く
		number, err := strconv.Atoi(input) // 入力を整数に変換
		if err != nil {
			if core.Developing && core.Debug {
				log.Printf("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
			continue
		}
		// 数値が1から4の間かどうかをチェック
		if number >= 1 && number <= 4 {
			switch number {
			case 1:
				state.Color = core.BLUE__
			case 2:
				state.Color = core.GREEN_
			case 3:
				state.Color = core.RED___
			case 4:
				state.Color = core.YELLOW
			}
			color := core.ColorName(state.Color)
			fmt.Printf("> 次の色は %v です。\n", color)
			break
		}
		if core.Developing && core.Debug {
			log.Printf("\t\t> 1から4の数を指定して下さい。")
		}
		fmt.Printf("\t\t> 1から4の数を指定して下さい。\n")
	}
}
