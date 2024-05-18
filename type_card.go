package main

import (
	"fmt"
	"log"

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
	WILD___ int = iota + 100 // 100 ワイルド         8枚
	DRAW4__                  // 101 ワイルドドロー4  4枚
	DRAW2__                  // 1*2 ドロー2    各2枚
	REVERSE                  // 1*3 リバース2  各2枚
	SKIP___                  // 1*4 スキップ2  各2枚
	DISCARD                  // 1*5 ディスカードオール (公式アプリGOワイルドルール) 正確な枚数不明。
)

const (
	effect_wild       = "色を指定できます。"
	effect_d4         = "次の人は4枚引き、スキップします。色を指定できます。"
	effect_d2         = "次の人は2枚引き、スキップします。"
	effect_reverse    = "順番の移る方向が逆転します。"
	effect_skip       = "次の人をスキップします。"
	effect_discardall = "手札にある同色のカードを全て捨てます。"
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
		return "無"
	}
	switch c {
	case BLUE__:
		return "青"
	case GREEN_:
		return "緑"
	case RED___:
		return "赤"
	case YELLOW:
		return "黄"
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

func (c *Card) ActionEffect(limit *State, hand *[]Card) {
	if developing && debug {
		fmt.Printf("\n")
		log.Printf("c.ActionType(): %v", c.ActionType())
	}
	if !c.IsAction() {
		panic(errors.Errorf("%v", "バグ"))
	}
	action := c.ActionType() + 100
	switch action {
	case WILD___:
		PromptWild(limit)
	case DRAW4__:
		PromptWild(limit)
		limit.Draw = 4
	case DRAW2__:
		limit.Color = c.Color()
		limit.Draw = 2
	case REVERSE:
		limit.Reverse()
	case SKIP___:
		limit.Skip()
	case DISCARD:
		limit.DiscardAll(c, hand)
	default:
		log.Printf("action:%v", action)
		panic(errors.Errorf("%v", "バグ"))
	}
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

func GetAllCards() ([]Card, []Card) {

	official := []Card{ // カードは112枚
		// スライスの要素番号をIDとする。
		// ワイルド // 8枚
		/*   0 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   1 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   2 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   3 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   4 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   5 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   6 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/*   7 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		// ドロー4 // 4枚
		/*   8 */ {"WD+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		/*   9 */ {"WD+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		/*  10 */ {"WD+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		/*  11 */ {"WD+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		// ドロー2 // 各色2枚
		/*  12 */ {"bD+2", DRAW2__ + BLUE__, 20, "青ドロー2", effect_d2},
		/*  13 */ {"bD+2", DRAW2__ + BLUE__, 20, "青ドロー2", effect_d2},
		/*  14 */ {"gD+2", DRAW2__ + GREEN_, 20, "緑ドロー2", effect_d2},
		/*  15 */ {"gD+2", DRAW2__ + GREEN_, 20, "緑ドロー2", effect_d2},
		/*  16 */ {"rD+2", DRAW2__ + RED___, 20, "赤ドロー2", effect_d2},
		/*  17 */ {"rD+2", DRAW2__ + RED___, 20, "赤ドロー2", effect_d2},
		/*  18 */ {"yD+2", DRAW2__ + YELLOW, 20, "黄ドロー2", effect_d2},
		/*  19 */ {"yD+2", DRAW2__ + YELLOW, 20, "黄ドロー2", effect_d2},
		// リバース // 各色2枚
		/*  20 */ {"bREV", REVERSE + BLUE__, 20, "青リバース", effect_reverse},
		/*  21 */ {"bREV", REVERSE + BLUE__, 20, "青リバース", effect_reverse},
		/*  22 */ {"gREV", REVERSE + GREEN_, 20, "緑リバース", effect_reverse},
		/*  23 */ {"gREV", REVERSE + GREEN_, 20, "緑リバース", effect_reverse},
		/*  24 */ {"rREV", REVERSE + RED___, 20, "赤リバース", effect_reverse},
		/*  25 */ {"rREV", REVERSE + RED___, 20, "赤リバース", effect_reverse},
		/*  26 */ {"yREV", REVERSE + YELLOW, 20, "黄リバース", effect_reverse},
		/*  27 */ {"yREV", REVERSE + YELLOW, 20, "黄リバース", effect_reverse},
		// スキップ // 各色2枚
		/*  28 */ {"bSKP", SKIP___ + BLUE__, 20, "青スキップ", effect_skip},
		/*  29 */ {"bSKP", SKIP___ + BLUE__, 20, "青スキップ", effect_skip},
		/*  30 */ {"gSKP", SKIP___ + GREEN_, 20, "緑スキップ", effect_skip},
		/*  31 */ {"gSKP", SKIP___ + GREEN_, 20, "緑スキップ", effect_skip},
		/*  32 */ {"rSKP", SKIP___ + RED___, 20, "赤スキップ", effect_skip},
		/*  33 */ {"rSKP", SKIP___ + RED___, 20, "赤スキップ", effect_skip},
		/*  34 */ {"ySKP", SKIP___ + YELLOW, 20, "黄スキップ", effect_skip},
		/*  35 */ {"ySKP", SKIP___ + YELLOW, 20, "黄スキップ", effect_skip},
		// 数字 各色19枚 0だけ1枚で他2枚
		/*  36 */ {"b0__", 0 + BLUE__, 0, "青0", ""},
		/*  37 */ {"b1__", 1 + BLUE__, 1, "青1", ""},
		/*  38 */ {"b1__", 1 + BLUE__, 1, "青1", ""},
		/*  39 */ {"b2__", 2 + BLUE__, 2, "青2", ""},
		/*  40 */ {"b2__", 2 + BLUE__, 2, "青2", ""},
		/*  41 */ {"b3__", 3 + BLUE__, 3, "青3", ""},
		/*  42 */ {"b3__", 3 + BLUE__, 3, "青3", ""},
		/*  43 */ {"b4__", 4 + BLUE__, 4, "青4", ""},
		/*  44 */ {"b4__", 4 + BLUE__, 4, "青4", ""},
		/*  45 */ {"b5__", 5 + BLUE__, 5, "青5", ""},
		/*  46 */ {"b5__", 5 + BLUE__, 5, "青5", ""},
		/*  47 */ {"b6__", 6 + BLUE__, 6, "青6", ""},
		/*  48 */ {"b6__", 6 + BLUE__, 6, "青6", ""},
		/*  49 */ {"b7__", 7 + BLUE__, 7, "青7", ""},
		/*  50 */ {"b7__", 7 + BLUE__, 7, "青7", ""},
		/*  51 */ {"b8__", 8 + BLUE__, 8, "青8", ""},
		/*  52 */ {"b8__", 8 + BLUE__, 8, "青8", ""},
		/*  53 */ {"b9__", 9 + BLUE__, 9, "青9", ""},
		/*  54 */ {"b9__", 9 + BLUE__, 9, "青9", ""},
		/*  55 */ {"g0__", 0 + GREEN_, 0, "緑0", ""},
		/*  56 */ {"g1__", 1 + GREEN_, 1, "緑1", ""},
		/*  57 */ {"g1__", 1 + GREEN_, 1, "緑1", ""},
		/*  58 */ {"g2__", 2 + GREEN_, 2, "緑2", ""},
		/*  59 */ {"g2__", 2 + GREEN_, 2, "緑2", ""},
		/*  60 */ {"g3__", 3 + GREEN_, 3, "緑3", ""},
		/*  61 */ {"g3__", 3 + GREEN_, 3, "緑3", ""},
		/*  62 */ {"g4__", 4 + GREEN_, 4, "緑4", ""},
		/*  63 */ {"g4__", 4 + GREEN_, 4, "緑4", ""},
		/*  64 */ {"g5__", 5 + GREEN_, 5, "緑5", ""},
		/*  65 */ {"g5__", 5 + GREEN_, 5, "緑5", ""},
		/*  66 */ {"g6__", 6 + GREEN_, 6, "緑6", ""},
		/*  67 */ {"g6__", 6 + GREEN_, 6, "緑6", ""},
		/*  68 */ {"g7__", 7 + GREEN_, 7, "緑7", ""},
		/*  69 */ {"g7__", 7 + GREEN_, 7, "緑7", ""},
		/*  70 */ {"g8__", 8 + GREEN_, 8, "緑8", ""},
		/*  71 */ {"g8__", 8 + GREEN_, 8, "緑8", ""},
		/*  72 */ {"g9__", 9 + GREEN_, 9, "緑9", ""},
		/*  73 */ {"g9__", 9 + GREEN_, 9, "緑9", ""},
		/*  74 */ {"r0__", 0 + RED___, 0, "赤0", ""},
		/*  75 */ {"r1__", 1 + RED___, 1, "赤1", ""},
		/*  76 */ {"r1__", 1 + RED___, 1, "赤1", ""},
		/*  77 */ {"r2__", 2 + RED___, 2, "赤2", ""},
		/*  78 */ {"r2__", 2 + RED___, 2, "赤2", ""},
		/*  79 */ {"r3__", 3 + RED___, 3, "赤3", ""},
		/*  80 */ {"r3__", 3 + RED___, 3, "赤3", ""},
		/*  81 */ {"r4__", 4 + RED___, 4, "赤4", ""},
		/*  82 */ {"r4__", 4 + RED___, 4, "赤4", ""},
		/*  83 */ {"r5__", 5 + RED___, 5, "赤5", ""},
		/*  84 */ {"r5__", 5 + RED___, 5, "赤5", ""},
		/*  85 */ {"r6__", 6 + RED___, 6, "赤6", ""},
		/*  86 */ {"r6__", 6 + RED___, 6, "赤6", ""},
		/*  87 */ {"r7__", 7 + RED___, 7, "赤7", ""},
		/*  88 */ {"r7__", 7 + RED___, 7, "赤7", ""},
		/*  89 */ {"r8__", 8 + RED___, 8, "赤8", ""},
		/*  90 */ {"r8__", 8 + RED___, 8, "赤8", ""},
		/*  91 */ {"r9__", 9 + RED___, 9, "赤9", ""},
		/*  92 */ {"r9__", 9 + RED___, 9, "赤9", ""},
		/*  93 */ {"y0__", 0 + YELLOW, 0, "黄0", ""},
		/*  94 */ {"y1__", 1 + YELLOW, 1, "黄1", ""},
		/*  95 */ {"y1__", 1 + YELLOW, 1, "黄1", ""},
		/*  96 */ {"y2__", 2 + YELLOW, 2, "黄2", ""},
		/*  97 */ {"y2__", 2 + YELLOW, 2, "黄2", ""},
		/*  98 */ {"y3__", 3 + YELLOW, 3, "黄3", ""},
		/*  99 */ {"y3__", 3 + YELLOW, 3, "黄3", ""},
		/* 100 */ {"y4__", 4 + YELLOW, 4, "黄4", ""},
		/* 101 */ {"y4__", 4 + YELLOW, 4, "黄4", ""},
		/* 102 */ {"y5__", 5 + YELLOW, 5, "黄5", ""},
		/* 103 */ {"y5__", 5 + YELLOW, 5, "黄5", ""},
		/* 104 */ {"y6__", 6 + YELLOW, 6, "黄6", ""},
		/* 105 */ {"y6__", 6 + YELLOW, 6, "黄6", ""},
		/* 106 */ {"y7__", 7 + YELLOW, 7, "黄7", ""},
		/* 107 */ {"y7__", 7 + YELLOW, 7, "黄7", ""},
		/* 108 */ {"y8__", 8 + YELLOW, 8, "黄8", ""},
		/* 109 */ {"y8__", 8 + YELLOW, 8, "黄8", ""},
		/* 110 */ {"y9__", 9 + YELLOW, 9, "黄9", ""},
		/* 111 */ {"y9__", 9 + YELLOW, 9, "黄9", ""},
		// 公式標準ルールで使用されるカードはここまで。
	}

	// 公式アプリ内のGOワイルドで使用されるカード
	// * 公式ワイルドルール
	// * 2倍のデッキ。
	// * +2に+2,+4を出せる。
	// * アクションカード枚数増加
	// * ディスカードの追加。同じ色のカードをすべて捨てられる。
	wild := make([]Card, 9, len(official)*3)
	wild = append(wild, official...) // 2倍のデッキ。 // IDは0から111まで
	wild = append(wild, official...) // 2倍のデッキ。 // IDは112から223まで
	// * アクションカード枚数増加
	// * ディスカードの追加。同じ色のカードをすべて捨てられる。各色2枚と思われる。
	// アクションカード枚数増加と書いてあるがどれだけ増えているのか不明。
	// 元の112枚中36(8+4+8+8+8)枚の32%がアクションカードなので、4割程度になるように増加させると想定する。
	// 2倍のデッキの224(112x2)枚中アクションカードは72(36x2)枚で、4割は90枚程度なので20(90-72)枚弱増やす場合、
	// 元のアクションカードの半分の18(4+2+4+4+4)が妥当と思われる。
	// また、これにディスカードをほかのアクションカード同じ12枚入れると想定すると、
	// 数字カード152(19x4x2)枚、アクションカード102(12+6+12+12+12+12)枚の合計254枚。
	// これで、102/152=0.401なのでちょうど4割になる。

	additional := []Card{ // IDは224から253まで
		// アクションカード増加
		// ワイルド // 4枚
		/* 224 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/* 225 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/* 226 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/* 227 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		// ドロー4 // 2枚
		/* 228 */ {"*D+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		/* 229 */ {"*D+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		// ドロー2 // 各色1枚
		/* 230 */ {"bD+2", DRAW2__ + BLUE__, 20, "青ドロー2", effect_d2},
		/* 231 */ {"gD+2", DRAW2__ + GREEN_, 20, "緑ドロー2", effect_d2},
		/* 232 */ {"rD+2", DRAW2__ + RED___, 20, "赤ドロー2", effect_d2},
		/* 233 */ {"yD+2", DRAW2__ + YELLOW, 20, "黄ドロー2", effect_d2},
		// リバース // 各色1枚
		/* 234 */ {"bREV", REVERSE + BLUE__, 20, "青リバース", effect_reverse},
		/* 235 */ {"gREV", REVERSE + GREEN_, 20, "緑リバース", effect_reverse},
		/* 236 */ {"rREV", REVERSE + RED___, 20, "赤リバース", effect_reverse},
		/* 237 */ {"yREV", REVERSE + YELLOW, 20, "黄リバース", effect_reverse},
		// スキップ // 各色1枚
		/* 238 */ {"bSKP", SKIP___ + BLUE__, 20, "青スキップ", effect_skip},
		/* 239 */ {"gSKP", SKIP___ + GREEN_, 20, "緑スキップ", effect_skip},
		/* 240 */ {"rSKP", SKIP___ + RED___, 20, "赤スキップ", effect_skip},
		/* 241 */ {"ySKP", SKIP___ + YELLOW, 20, "黄スキップ", effect_skip},
		// ディスカード増加 各色3枚
		/* 242 */ {"bDIS", DISCARD + BLUE__, 20, "青ディスカードオール", effect_discardall},
		/* 243 */ {"bDIS", DISCARD + BLUE__, 20, "青ディスカードオール", effect_discardall},
		/* 244 */ {"bDIS", DISCARD + BLUE__, 20, "青ディスカードオール", effect_discardall},
		/* 245 */ {"gDIS", DISCARD + GREEN_, 20, "緑ディスカードオール", effect_discardall},
		/* 246 */ {"gDIS", DISCARD + GREEN_, 20, "緑ディスカードオール", effect_discardall},
		/* 247 */ {"gDIS", DISCARD + GREEN_, 20, "緑ディスカードオール", effect_discardall},
		/* 248 */ {"rDIS", DISCARD + RED___, 20, "赤ディスカードオール", effect_discardall},
		/* 249 */ {"rDIS", DISCARD + RED___, 20, "赤ディスカードオール", effect_discardall},
		/* 250 */ {"rDIS", DISCARD + RED___, 20, "赤ディスカードオール", effect_discardall},
		/* 251 */ {"yDIS", DISCARD + YELLOW, 20, "黄ディスカードオール", effect_discardall},
		/* 252 */ {"yDIS", DISCARD + YELLOW, 20, "黄ディスカードオール", effect_discardall},
		/* 253 */ {"yDIS", DISCARD + YELLOW, 20, "黄ディスカードオール", effect_discardall},
	}
	wild = append(wild, additional...)
	return official, wild
}
