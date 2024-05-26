package core

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

// 色が決まらないカードか。
// ワイルドドロー4、シャッフルワイルド、ホワイトワイルド
func (c *Card) IsWild() bool {
	if c.Type >= 100 && c.Type < 200 { // 100以上200未満の時。
		return c.IsActionCardType([]int{0, 1, 6, 8, 9}, 100, 200)
	}
	// 200以上の時。カードが未定義。定義されてたら実装する。
	return false
}

func (c *Card) IsSkip() bool {
	return c.IsActionCardType([]int{4}, 100, 200)
}
func (c *Card) IsReverse() bool {
	return c.IsActionCardType([]int{3}, 100, 200)
}
func (c *Card) IsDraw() bool {
	return c.IsActionCardType([]int{1, 2, 7, 8, 9}, 100, 200)
}
func (c *Card) IsDiscardAll() bool {
	return c.IsActionCardType([]int{5}, 100, 200)
}
func (c *Card) IsShuffleWild() bool {
	return c.IsActionCardType([]int{6}, 100, 200)
}
func (c *Card) IsActionCardType(typenum []int, min_num int, max_num int) bool {
	if c.IsAction() {
		if min_num <= c.Type && c.Type < max_num {
			t := c.ActionType()
			for _, n := range typenum {
				if t == n {
					return true
				}
			}
		}
	}
	return false
}
