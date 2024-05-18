package main

// push関数：任意の型の要素を受け取るためにCardを使用
func push(slice *[]Card, value Card) {
	*slice = append(*slice, value)
}

// pop関数：スライスから要素を削除
func pop(slice *[]Card) (Card, bool) {
	if len(*slice) == 0 {
		return Card{}, false // スライスが空の場合は何もしない
	}
	// スライスの最後の要素を取得し、スライスからその要素を除去
	value := (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
	return value, true
}

// peek関数：スタックの一番上を覗く
func peek(slice *[]Card) (Card, bool) {
	if len(*slice) == 0 {
		return Card{}, false // スライスが空の場合は何もしない
	}
	return (*slice)[len(*slice)-1], true
}
