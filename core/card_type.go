package core

import (
	"log"

	"github.com/pkg/errors"
)

type CardInfo struct {
	CardType      []Card
	CardNameList  []string
	CardTypeCount map[string]int // カードの種類と枚数
}

func (ci *CardInfo) ID(name string) int {
	for i, card := range ci.CardType {
		if name == card.Name {
			return i
		}
	}
	panic(errors.Errorf("%v", "バグ")) // 無かったらエラー
	return -1
}
func (ci *CardInfo) Get(name string) *Card {
	id := ci.ID(name)
	if id < 0 {
		panic(errors.Errorf("%v", "バグ")) // 無かったらエラー
	}
	return &ci.CardType[id]
}
func (ci *CardInfo) AddCardType(card Card, num int) {
	if id := ci.ID(card.Name); id >= 0 {
		log.Print("name", card.Name)
		panic(errors.Errorf("%v", "既存")) // すでにあったらエラー
	}
	ci.CardType = append(ci.CardType, card)
	ci.CardNameList = append(ci.CardNameList, card.Name)
	ci.CardTypeCount[card.Name] = num
}
func (ci *CardInfo) AddNumberOfCard(name string, num int) {
	if id := ci.ID(name); id < 0 {
		panic(errors.Errorf("%v", "バグ")) // 無かったらエラー
	}
	ci.CardTypeCount[name] += num
}
func (ci *CardInfo) SetNumberOfCard(name string, num int) {
	if id := ci.ID(name); id < 0 {
		panic(errors.Errorf("%v", "バグ")) // 無かったらエラー
	}
	ci.CardTypeCount[name] = num
}
func (ci *CardInfo) GenerateDeck() []Card {
	deck := []Card{}
	for _, name := range ci.CardNameList {
		for i := 0; i < ci.CardTypeCount[name]; i++ {
			deck = append(deck, *ci.Get(name))
		}
	}
	return deck
}

func (ci *CardInfo) Init() {
	ci.CardType = []Card{ // カードは112枚
		/*  0 */ {"b0__", 0 + BLUE__, 0, "青0", ""}, // 数字 各色19枚 0だけ1枚で他2枚
		/*  1 */ {"b1__", 1 + BLUE__, 1, "青1", ""},
		/*  2 */ {"b2__", 2 + BLUE__, 2, "青2", ""},
		/*  3 */ {"b3__", 3 + BLUE__, 3, "青3", ""},
		/*  4 */ {"b4__", 4 + BLUE__, 4, "青4", ""},
		/*  5 */ {"b5__", 5 + BLUE__, 5, "青5", ""},
		/*  6 */ {"b6__", 6 + BLUE__, 6, "青6", ""},
		/*  7 */ {"b7__", 7 + BLUE__, 7, "青7", ""},
		/*  8 */ {"b8__", 8 + BLUE__, 8, "青8", ""},
		/*  9 */ {"b9__", 9 + BLUE__, 9, "青9", ""},
		/* 10 */ {"g0__", 0 + GREEN_, 0, "緑0", ""},
		/* 11 */ {"g1__", 1 + GREEN_, 1, "緑1", ""},
		/* 12 */ {"g2__", 2 + GREEN_, 2, "緑2", ""},
		/* 13 */ {"g3__", 3 + GREEN_, 3, "緑3", ""},
		/* 14 */ {"g4__", 4 + GREEN_, 4, "緑4", ""},
		/* 15 */ {"g5__", 5 + GREEN_, 5, "緑5", ""},
		/* 16 */ {"g6__", 6 + GREEN_, 6, "緑6", ""},
		/* 17 */ {"g7__", 7 + GREEN_, 7, "緑7", ""},
		/* 18 */ {"g8__", 8 + GREEN_, 8, "緑8", ""},
		/* 19 */ {"g9__", 9 + GREEN_, 9, "緑9", ""},
		/* 20 */ {"r0__", 0 + RED___, 0, "赤0", ""},
		/* 21 */ {"r1__", 1 + RED___, 1, "赤1", ""},
		/* 22 */ {"r2__", 2 + RED___, 2, "赤2", ""},
		/* 23 */ {"r3__", 3 + RED___, 3, "赤3", ""},
		/* 24 */ {"r4__", 4 + RED___, 4, "赤4", ""},
		/* 25 */ {"r5__", 5 + RED___, 5, "赤5", ""},
		/* 26 */ {"r6__", 6 + RED___, 6, "赤6", ""},
		/* 27 */ {"r7__", 7 + RED___, 7, "赤7", ""},
		/* 28 */ {"r8__", 8 + RED___, 8, "赤8", ""},
		/* 29 */ {"r9__", 9 + RED___, 9, "赤9", ""},
		/* 30 */ {"y0__", 0 + YELLOW, 0, "黄0", ""},
		/* 31 */ {"y1__", 1 + YELLOW, 1, "黄1", ""},
		/* 32 */ {"y2__", 2 + YELLOW, 2, "黄2", ""},
		/* 33 */ {"y3__", 3 + YELLOW, 3, "黄3", ""},
		/* 34 */ {"y4__", 4 + YELLOW, 4, "黄4", ""},
		/* 35 */ {"y5__", 5 + YELLOW, 5, "黄5", ""},
		/* 36 */ {"y6__", 6 + YELLOW, 6, "黄6", ""},
		/* 37 */ {"y7__", 7 + YELLOW, 7, "黄7", ""},
		/* 38 */ {"y8__", 8 + YELLOW, 8, "黄8", ""},
		/* 39 */ {"y9__", 9 + YELLOW, 9, "黄9", ""},
		/* 40 */ {"WILD", WILD___, 50, "ワイルドカード", effect_wild},
		/* 41 */ {"WD+4", DRAW4__, 50, "ワイルドドロー4", effect_d4},
		/* 42 */ {"bD+2", DRAW2__ + BLUE__, 20, "青ドロー2", effect_d2},
		/* 43 */ {"bREV", REVERSE + BLUE__, 20, "青リバース", effect_reverse},
		/* 44 */ {"bSKP", SKIP___ + BLUE__, 20, "青スキップ", effect_skip},
		/* 45 */ {"gD+2", DRAW2__ + GREEN_, 20, "緑ドロー2", effect_d2},
		/* 46 */ {"gREV", REVERSE + GREEN_, 20, "緑リバース", effect_reverse},
		/* 47 */ {"gSKP", SKIP___ + GREEN_, 20, "緑スキップ", effect_skip},
		/* 48 */ {"rD+2", DRAW2__ + RED___, 20, "赤ドロー2", effect_d2},
		/* 49 */ {"rREV", REVERSE + RED___, 20, "赤リバース", effect_reverse},
		/* 50 */ {"rSKP", SKIP___ + RED___, 20, "赤スキップ", effect_skip},
		/* 51 */ {"yD+2", DRAW2__ + YELLOW, 20, "黄ドロー2", effect_d2},
		/* 52 */ {"yREV", REVERSE + YELLOW, 20, "黄リバース", effect_reverse},
		/* 53 */ {"ySKP", SKIP___ + YELLOW, 20, "黄スキップ", effect_skip},
		/* 54 */ {"bDIS", DISCARD + BLUE__, 20, "青ディスカードオール", effect_discardall},
		/* 55 */ {"gDIS", DISCARD + GREEN_, 20, "緑ディスカードオール", effect_discardall},
		/* 56 */ {"rDIS", DISCARD + RED___, 20, "赤ディスカードオール", effect_discardall},
		/* 57 */ {"yDIS", DISCARD + YELLOW, 20, "黄ディスカードオール", effect_discardall},
		/* 58 */ {"Shff", SHUFFLE, 50, "シャッフルワイルド", effect_shuffle_wild},
		/* 59 */ {"bHT2", HIT2___ + BLUE__, 20, "青ヒット2", effect_hit_2},
		/* 60 */ {"gHT2", HIT2___ + GREEN_, 20, "緑ヒット2", effect_hit_2},
		/* 61 */ {"rHT2", HIT2___ + RED___, 20, "赤ヒット2", effect_hit_2},
		/* 62 */ {"yHT2", HIT2___ + YELLOW, 20, "黄ヒット2", effect_hit_2},
		/* 63 */ {"WHT4", WILD_H4, 50, "ワイルドヒット4", effect_wild_hit_4},
		/* 64 */ {"WAA2", WILD_AA, 50, "ワイルド アタックアタックカード", effect_wild_attack},
	}
	ci.CardTypeCount = map[string]int{
		"b0__": 1, "b1__": 2, "b2__": 2, "b3__": 2, "b4__": 2, "b5__": 2, "b6__": 2, "b7__": 2, "b8__": 2, "b9__": 2,
		"g0__": 1, "g1__": 2, "g2__": 2, "g3__": 2, "g4__": 2, "g5__": 2, "g6__": 2, "g7__": 2, "g8__": 2, "g9__": 2,
		"r0__": 1, "r1__": 2, "r2__": 2, "r3__": 2, "r4__": 2, "r5__": 2, "r6__": 2, "r7__": 2, "r8__": 2, "r9__": 2,
		"y0__": 1, "y1__": 2, "y2__": 2, "y3__": 2, "y4__": 2, "y5__": 2, "y6__": 2, "y7__": 2, "y8__": 2, "y9__": 2,
		"WILD": 8, "WD+4": 4,
		"bD+2": 2, "bREV": 2, "bSKP": 2,
		"gD+2": 2, "gREV": 2, "gSKP": 2,
		"rD+2": 2, "rREV": 2, "rSKP": 2,
		"yD+2": 2, "yREV": 2, "ySKP": 2,
	}
	ci.CardNameList = []string{
		"b0__", "b1__", "b2__", "b3__", "b4__", "b5__", "b6__", "b7__", "b8__", "b9__",
		"g0__", "g1__", "g2__", "g3__", "g4__", "g5__", "g6__", "g7__", "g8__", "g9__",
		"r0__", "r1__", "r2__", "r3__", "r4__", "r5__", "r6__", "r7__", "r8__", "r9__",
		"y0__", "y1__", "y2__", "y3__", "y4__", "y5__", "y6__", "y7__", "y8__", "y9__",
		"WILD", "WD+4",
		"bD+2", "bREV", "bSKP",
		"gD+2", "gREV", "gSKP",
		"rD+2", "rREV", "rSKP",
		"yD+2", "yREV", "ySKP",
	}
}
