package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func promptNumber() int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a number (1-4): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		// 改行文字を取り除く
		input = strings.TrimSpace(input)

		// 入力を整数に変換
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		// 数値が1から4の間かどうかをチェック
		if number >= 1 && number <= 4 {
			return number
		}

		fmt.Println("Please enter a number between 1 and 4.")
	}
}

func PromptWild(limit *State) {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\t\t> 色を指定してください。:\n")
		fmt.Printf("\t\t> 青:1 / 緑:2 / 赤:3 / 黄:4 > ")

		var input string
		var err error

		if debug {
			// デバッグ中はとりあえずランダムに選択。
			input = fmt.Sprintf("%v", rand.Intn(5)) // 0-4のはず。 0の時再度プロンプト表示。
			fmt.Printf("%v\n", input)
		} else {
			input, err = reader.ReadString('\n')
		}

		if err != nil {
			if developing && debug {
				log.Println("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Println("\t\t> 1から4の数を指定して下さい。")
			continue
		}
		input = strings.TrimSpace(input)   // 改行文字を取り除く
		number, err := strconv.Atoi(input) // 入力を整数に変換
		if err != nil {
			if developing && debug {
				log.Println("\t\t> 1から4の数を指定して下さい。")
			}
			fmt.Println("\t\t> 1から4の数を指定して下さい。")
			continue
		}
		// 数値が1から4の間かどうかをチェック
		if number >= 1 && number <= 4 {
			switch number {
			case 1:
				limit.Color = BLUE__
			case 2:
				limit.Color = GREEN_
			case 3:
				limit.Color = RED___
			case 4:
				limit.Color = YELLOW
			}
			color := ColorName(limit.Color)
			fmt.Printf("\t\t> 次の色は %v です。\n", color)
			break
		}
		if developing && debug {
			log.Println("\t\t> 1から4の数を指定して下さい。")
		}
		fmt.Println("\t\t> 1から4の数を指定して下さい。")
	}
}
