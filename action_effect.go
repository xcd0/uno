package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func ActionEffect(c *Card, state *State, npc bool) {
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
		SelectColor(state, npc) // 色を選択させる。
		state.Draw = 0          // 次の人は追加でカードを引く必要はない。
		state.Number = -1       // 次の人は数字を無視してよい。
	case DRAW4__:
		SelectColor(state, npc) // 色を選択させる。
		state.Draw += 4         // 次の人は追加でカードを4枚引く必要がある。
		state.Number = -1       // 次の人は数字を無視してよい。
	case DRAW2__:
		state.Color = c.Color() // 色はカードの色。
		state.Draw += 2         // 次の人は追加でカードを2枚引く必要がある。
		state.Number = -1       // 次の人は数字を無視してよい。
	case REVERSE:
		state.Reverse()
		state.Color = c.Color() // 色はカードの色。
		state.Draw = 0          // 次の人は追加でカードを引く必要はない。
		state.Number = -1       // 次の人は数字を無視してよい。
	case SKIP___:
		state.Skip()
		state.Color = c.Color() // 色はカードの色。
		state.Draw = 0          // 次の人は追加でカードを引く必要はない。
		state.Number = -1       // 次の人は数字を無視してよい。
	case DISCARD:
		state.DiscardAll(c)
		state.Color = c.Color() // 色はカードの色。
		state.Draw = 0          // 次の人は追加でカードを引く必要はない。
		state.Number = -1       // 次の人は数字を無視してよい。
	default:
		log.Printf("action: %v", action)
		panic(errors.Errorf("%v", "バグ"))
	}
}
