package core

import (
	crypto_rand "crypto/rand"
	"log"
	"math/big"
	math_rand "math/rand"
)

func ShuffleCards(cards []Card) {
	math_rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	if Developing && Debug {
		log.Printf("Shuffled Cards : len(cards): %v", len(cards))
		if len(cards) == 0 {
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
	if Developing && Debug {
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
