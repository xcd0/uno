package main

import (
	crypto_rand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"sort"

	"github.com/pkg/errors"
)

// 色
const (
	_NULL  = iota * 10
	BLUE__ // 10
	GREEN_ // 20
	RED___ // 30
	YELLOW // 40
)

// 種別
const (
	// 公式ルール
	WILD___ int = iota + 100 // 100 ワイルド
	DRAW4__                  // 101 ワイルドドロー4
	DRAW2__                  // 1*2 ドロー2
	REVERSE                  // 1*3 リバース2
	SKIP___                  // 1*4 スキップ2
	DISCARD                  // 1*5 ディスカードオール (公式アタックエクストリームルール,公式アプリGOワイルドルール)
	SHUFFLE                  // 1*6 シャッフルワイルド (公式追加ルール)
	HIT2___                  // 1*7 ヒット2 (公式アタックエクストリームルール)
	WILD_H4                  // 1*8 ワイルドヒット4 (公式アタックエクストリームルール)
	WILD_AA                  // 1*9 ワイルドアタックアタック (公式アタックエクストリームルール)
)

const (
	effect_wild         = "色を指定できます。"
	effect_d4           = "色を指定できます。次の人は4枚引き、スキップします。"
	effect_d2           = "次の人は2枚引き、スキップします。"
	effect_reverse      = "順番の移る方向が逆転します。"
	effect_skip         = "次の人をスキップします。プレーヤーが2人の時はスキップと同じになります。"
	effect_discardall   = "手札にある同色のカードを全て捨てます。"
	effect_shuffle_wild = "色を指定します。全員の手札をシャッフルして自分の左隣から順番に全て配ります。"
	effect_hit_2        = "次の人は2回ランダムな回数引き、スキップします。"
	effect_wild_hit_4   = "色を指定できます。次の人は4回ランダムな回数引き、スキップします。"
	effect_wild_attack  = "色を指定できます。指定した人は2回ランダムな回数引きます。順番は指定された次の人に移ります。"
)

// 上記を使用してカードの種類を整数で表現する
// 11 : 青1
// 12 : 青2
// 23 : 緑3
// 34 : 赤4
// 45 : 黄5
// 100 : ワイルド
// 101 : ワイルドドロー4
// 112 : 青ドロー2
// 123 : 緑リバース
// 134 : 赤スキップ

// したがって、下記のように判定できる。
// - 100未満は数字カード、100以上はアクションカード
// - 100未満の数字カードは、10の剰余を取って数値を判別する。
// - 100未満の数字カードは、intを10で割って10の剰余を取って色を判別する。
// - 100以上のアクションカードは、10の剰余を取って種類を判別する。

func ColorName(c int) string {
	if c < 0 {
		return "無色"
	}
	switch c {
	case BLUE__:
		return "青色"
	case GREEN_:
		return "緑色"
	case RED___:
		return "赤色"
	case YELLOW:
		return "黄色"
	}
	log.Printf("c:%v", c)
	panic(errors.Errorf("%v", "バグ"))
}

type Card struct {
	// ID   int    // 固有番号 カード全ての情報を保持するスライスの要素番号をそのままIDとして使用する。
	Name string // 表示するときの名前
	// カードの情報
	Type        int    // 種類
	Point       int    // 点数
	Description string // カードの種類のわかりやすい表示
	Detail      string // 説明
}

func Sort(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Name < cards[j].Name
	})

}

// 100未満は数字カード、100以上はアクションカード
func (c *Card) IsNum() bool {
	//if debug {
	//	log.Println("card: %v", *c)
	//}
	return c.Type < 100
}
func (c *Card) IsAction() bool {
	return !c.IsNum()
}

// 100未満の数字カードは、10の剰余を取って数値を判別する。
func (c *Card) Num() int {
	if !c.IsNum() {
		panic(errors.Errorf("%v", "バグ"))
	}
	return c.Type % 10
}

// 100未満の数字カードは、intを10で割って10の剰余を取って色を判別する。
func (c *Card) Color() int {
	// 数字か、色付きアクションカード
	if c.IsNum() {
	} else if c.Type%10 == DRAW2__%10 { // 1*2 ドロー2    各2枚
	} else if c.Type%10 == REVERSE%10 { // 1*3 リバース2  各2枚
	} else if c.Type%10 == SKIP___%10 { // 1*4 スキップ2  各2枚
	} else if c.Type%10 == DISCARD%10 { // 1*5 ディスカードオール (公式アプリGOワイルドルール) 正確な枚数不明。
	} else {
		log.Printf("card: %v", *c)
		panic(errors.Errorf("%v", "バグ"))
	}
	return ((c.Type / 10) % 10) * 10
}
func (c *Card) ColorName() string {
	return ColorName(c.Color())
}

// 100以上のアクションカードは、10の剰余を取って種類を判別する。
func (c *Card) ActionType() int {
	if !c.IsAction() {
		panic(errors.Errorf("%v", "バグ"))
	}
	return (c.Type / 10) % 10
}

func (c *Card) Print() string {
	return fmt.Sprintf("%v", c.Name)
}

func PrintCards(cards []Card) string {
	s := ""
	for i, c := range cards {
		s += c.Print()
		if i != len(cards)-1 {
			s += ","
		}
	}
	return s
}
func PrintCardsDescription(cards []Card) string {
	s := ""
	for i, c := range cards {
		s += c.Description
		if i != len(cards)-1 {
			s += ","
		}
	}
	return s
}

func PrintTopCards(cards []Card, num int) string {
	s := ""
	for i := len(cards) - 1; i >= len(cards)-num; i-- {
		c := cards[i]
		s += c.Print()
		if i != len(cards)-num {
			s += ","
		}
	}
	return s
}

//// RemoveCard は与えられたカードをスライスから削除します。
//func RemoveCard(cards *[]Card, target Card) {
//	log.Printf("RemoveCard Pre : %v", PrintCards(*cards))
//	for i, card := range *cards {
//		if card.Name == target.Name { // カードがマッチしたら、その位置でスライスを分割し、該当カードを除外
//			*cards = append((*cards)[:i], (*cards)[i+1:]...)
//		}
//	}
//	log.Printf("RemoveCard Post: %v", PrintCards(*cards))
//}

// RemoveCard は与えられたカードをスライスから削除します。
func RemoveCard(cards []Card, target Card) []Card {
	for i, card := range cards {
		if card.Name == target.Name { // カードがマッチしたら、その位置でスライスを分割し、該当カードを除外
			// カードがマッチしたら、その位置でスライスを分割し、該当カードを除外
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards // 該当するカードがなければ元のスライスを返す
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
