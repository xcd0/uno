package main

import (
	"fmt"
)

type State struct {
	Color          int      // -1の時任意
	Number         int      // -1の時任意
	Draw           int      // 0以下の時任意
	Turn           int      // CW:1 CCW:-1
	NowPlayerIndex int      // 今の人のID
	Players        []string // 参加者
	Hands          [][]Card // 手札
	Deck           []Card   // 山札
	Discard        []Card   // 捨て札
}

func (l *State) Print() string {
	num := ""
	if l.Number >= 0 {
		num = fmt.Sprintf("%v", l.Number)
	}
	turn := "CW"
	if l.Turn < 0 {
		turn = "CCW"
	}
	return fmt.Sprintf("(%v)%v%v+%v%v",
		l.Players[l.NowPlayerIndex], // 今の人のID
		ColorName(l.Color),          // -1の時任意
		num,                         // -1の時任意
		l.Draw,                      // 0以下の時任意
		turn,                        // CW:1 CCW:-1
	)
}

func (l *State) Reverse() {
	fmt.Printf("\t> 順番が反転します。\n")
	if l.Turn > 0 {
		l.Turn = -1
	} else {
		l.Turn = 1
	}
}
func (l *State) Skip() {
	fmt.Printf("\t> %v さんはスキップされます。\n", l.NextName())
	l.NowPlayerIndex = ((l.NowPlayerIndex + l.Turn*2) + len(l.Players)*2) % len(l.Players)
}
func (l *State) DiscardAll(c *Card, hand *[]Card) {
	fmt.Printf("\t> 手札の %v のカード全てを捨てます。\n", c.ColorName())
}
func (l *State) Update(state *State) {
	l.NowPlayerIndex = ((l.NowPlayerIndex + l.Turn) + len(l.Players)) % len(l.Players)
	if l.Draw > 0 {
		fmt.Printf("\t> %v さんは%v枚引きます。\n", l.NowPlayerIndex, l.Draw)
		DealCards(l.NowPlayerIndex, l.Draw, state)
		l.Draw = 0
		l.Skip()
	}
}
func (l *State) Name() string {
	return l.Players[l.NowPlayerIndex]
}
func (l *State) NextName() string {
	tmp := ((l.NowPlayerIndex + l.Turn) + len(l.Players)) % len(l.Players)
	return l.Players[tmp]
}
