package main

import (
	crypto_rand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	math_rand "math/rand"
	"time"

	"github.com/pkg/errors"
)

var debug = true
var developing = false

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	math_rand.Seed(time.Now().UnixNano())

	official, wild := GetAllCards()

	players := []string{"You", "P1", "P2", "P3"}
	if true {
		GameStart(official, players)
	} else {
		GameStart(wild, players)
	}
}

func ShuffleCards(cards []Card) {
	math_rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	if developing && debug {
		if len(cards) == 0 {
			log.Printf("Shuffled Cards : len(cards): %v", len(cards))
			return
		}
		l := len(cards) / 8
		log.Printf("Shuffled Cards : \n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v",
			PrintCards(cards[:l]),
			PrintCards(cards[l+1:l*2]),
			PrintCards(cards[l*2+1:l*3]),
			PrintCards(cards[l*3+1:l*4]),
			PrintCards(cards[l*4+1:l*5]),
			PrintCards(cards[l*5+1:l*6]),
			PrintCards(cards[l*6+1:l*7]),
			PrintCards(cards[l*7+1:]),
		)
	}
}

// CryptoRandShuffle はカードのスライスを暗号学的に安全な方法でシャッフルします。
func CryptoRandShuffle(cards []Card) {
	n := len(cards)
	for i := range cards {
		// crypto/randを使用して、iからn-1の範囲の安全なランダムなインデックスを選びます
		jBig, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(n-i)))
		if err != nil {
			panic(err) // 乱数生成に失敗した場合
		}
		j := int(jBig.Int64()) + i

		// カードを交換
		cards[i], cards[j] = cards[j], cards[i]
	}
	if developing && debug {
		if len(cards) == 0 {
			log.Printf("CryptoRandShuffle: len(cards): %v", len(cards))
			return
		}
		l := len(cards) / 8
		log.Printf("CryptoRandShuffle: \n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v\n\t%#v",
			PrintCards(cards[:l]),
			PrintCards(cards[l+1:l*2]),
			PrintCards(cards[l*2+1:l*3]),
			PrintCards(cards[l*3+1:l*4]),
			PrintCards(cards[l*4+1:l*5]),
			PrintCards(cards[l*5+1:l*6]),
			PrintCards(cards[l*6+1:l*7]),
			PrintCards(cards[l*7+1:]),
		)
	}
}

// カードを引く。
func DealCards(player_id int, deal_num int, state *State) {

	// deal_num回state.Deckからhandにカードを引く。
	for j := 0; j < deal_num; j++ {
		if len(state.Deck) == 0 {
			// 捨て札から場にシャッフルして戻す。
			last_card, ok := pop(&state.Discard) // 捨て札の最も新しいカードを保持する。
			if !ok {
				panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???")) // 山札も捨て札もない。全てのカードが手札にある???
				// 引けないので引かないというのもありかなと思われる。
			}
			ShuffleCards(state.Discard) // 捨て札をシャッフル。
			CryptoRandShuffle(state.Discard)
			state.Deck = append(state.Deck, state.Discard...) // 山札にする。
			state.Discard = []Card{last_card}                 // 捨て札に保持しておいた最後の捨て札を戻す。
		}
		// 山札からカードを1枚引く。
		card, ok := pop(&state.Deck)
		push(&state.Hands[player_id], card)
		if !ok {
			// 謎。バグ。
			panic(errors.Errorf("%v", "バグ"))
		}
	}
	if developing && debug {
		log.Printf("Dealed Cards : %#v", PrintCards(state.Hands[player_id]))
	}
}

// InsertRandomCard はスライス cards のランダムな位置に新しいカード card を挿入します。
func InsertRandomCard(cards *[]Card, card Card) {
	n := len(*cards)
	// ランダムな位置を選択するためのインデックスを生成
	indexBig, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(n+1)))
	if err != nil {
		panic(err) // 乱数生成に失敗した場合はパニック
	}
	index := int(indexBig.Int64())

	// スライスにカードを挿入
	*cards = append((*cards)[:index], append([]Card{card}, (*cards)[index:]...)...)
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

func GameStart(cards []Card, players []string) {

	if developing && debug {
		log.Printf("Players Name: %#v", players)
		//log.Printf("Initial Cards: %#v", PrintCards(cards))
	}

	state := &State{
		Color:          -1,                           // -1の時任意
		Number:         -1,                           // -1の時任意
		Draw:           0,                            // 0以下の時任意
		Turn:           1,                            // CW:1 CCW:-1
		NowPlayerIndex: rand.Intn(len(players)),      // 今の人のID
		Players:        players,                      // 参加者
		Hands:          make([][]Card, len(players)), // 手札
		Deck:           make([]Card, len(cards)),     // 山札
		Discard:        make([]Card, 0, len(cards)),  // 捨て札
	}

	copy(state.Deck, cards)

	// shuffle
	ShuffleCards(state.Deck)
	CryptoRandShuffle(state.Deck)

	// 場からプレイヤーの手札を配る。初期7枚。
	for i := 0; i < len(players); i++ {
		DealCards(i, 7, state)
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
		last_card, ok := pop(&state.Deck)
		if !ok {
			panic(errors.Errorf("%v", "バグ"))
		}
		if last_card.IsAction() {
			InsertRandomCard(&state.Deck, last_card)
		} else {
			push(&state.Discard, last_card)
			break
		}
	}
	last_card, ok := peek(&state.Discard)
	if !ok {
		panic(errors.Errorf("%v", "バグ"))
	}

	// 5. 親の左隣の人が最初のプレイヤーです。時計周りにカードを捨てていきます。
	fmt.Printf("> %v さんが親です。親の次の人が最初です。\n", state.Name())
	fmt.Printf("> 最初のカードは %v です。\n", last_card.Name)
	state.Update(state)
	fmt.Printf("> %v さんのターンです。\n", state.Name())

	//for i := 0; i < 3; i++ {
	for {
		if debug {
			fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), PrintCards(state.Hands[state.NowPlayerIndex]))
		}

		// 捨てるカードの選択
		for {
			fmt.Printf("\t\t> 捨てるカードを選択してください。\n")

			if true {
				// とりあえず引かないでルール無視で1枚ずつ捨てるだけの実装。
				last_card, ok = pop(&state.Hands[state.NowPlayerIndex])
				if !ok {
					panic(errors.Errorf("%v", "バグ"))
				}
				push(&state.Discard, last_card)
				last_card, ok = peek(&state.Discard)
				if !ok {
					panic(errors.Errorf("%v", "バグ"))
				}
			} else {
				// ランダムに1枚選ぶ実装
				last_card, ok = PopRandomCard(&state.Hands[state.NowPlayerIndex])
				if !ok {
					panic(errors.Errorf("%v", "バグ"))
				}
			}

			// 選んだカードのチェック
			{
				// 未実装
				if false {
					// 再度選ぶ
					continue
				}
			}

			fmt.Printf("\t\t> %v を選択しました。\n", last_card.Print())
			break
		}

		// 捨てたカードがアクションカードならアクションカードの効果だけ発動させる。

		if last_card.IsNum() {
			//if developing && debug {
			//	log.Println("数値カード")
			//}
			state.Color = last_card.Color()
			state.Number = last_card.Num()
		} else if last_card.IsAction() {
			//if developing && debug {
			//	log.Println("アクションカード")
			//}
			//fmt.Printf("\t> アクションカード!\n")
			last_card.ActionEffect(state, &state.Hands[state.NowPlayerIndex])
		}

		{
			if len(state.Hands[state.NowPlayerIndex]) == 1 {
				fmt.Printf("\t> %v さんは %v を捨てました。UNO! です。\n", state.Name(), last_card.Print())
			} else {
				fmt.Printf("\t> %v さんは %v を捨てました。\n", state.Name(), last_card.Print())
			}

			if len(state.Hands[state.NowPlayerIndex]) == 0 {
				fmt.Printf("\t> %v さんの手札が無くなりました。%v さんの勝ちです。\n", state.Name(), state.Name())
				break
			} else if len(state.Hands[state.NowPlayerIndex]) == 1 {
				//fmt.Printf("\t> %v さんは UNO! です。\n", state.Name())
			} else {
				fmt.Printf("\t> %v さんの残り枚数は %v です。\n", state.Name(), len(state.Hands[state.NowPlayerIndex]))
			}
		}

		if debug {
			fmt.Printf("\t> %v さんの手札 : %v\n", state.Name(), PrintCards(state.Hands[state.NowPlayerIndex]))
		}
		fmt.Printf("\t---------------------------------> %v\n", state.Print())
		state.Update(state)
		fmt.Printf("> %v さんのターンです。\n", state.Name())
	}
}
