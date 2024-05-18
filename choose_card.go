package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func ChooseCard(state *State, npc bool) *Card {
	var chosen Card
	tmp := FilterCards(state)
	// forの間一時的にコピーしたものから取り出す。
	dealed := false
	for {
		if len(tmp) == 0 {
			if !dealed {
				dealed = true
				// 出せるカードがないので追加で1枚引く。
				fmt.Printf("\t> %v さんは出せるカードがないので1枚引きます。\n", state.Name())
				ok := DealCards(state.NowPlayerIndex, 1, state)
				if !ok {
					fmt.Printf("\t> 山札がないので %v さんはカードを引くことができませんでした。スキップします。\n", state.Name())
					return nil
				}
				continue
			} else {
				return nil
			}
		}
		if developing && debug {
			fmt.Printf("\t> %v さんの手札で出せるもの : %v\n", state.Name(), PrintCardsDescription(tmp))
		}

		if selected := SelectCard(state, tmp, npc); selected == nil {
			return nil
		} else {
			chosen = *selected
		}

		if CheckCard(chosen, state) {
			if !npc || developing && debug {
				fmt.Printf(" > OK!\n")
			}
			break
		}
		if !npc || developing && debug {
			fmt.Printf(" > NG!\n")
		}

	}
	*state.Hand() = RemoveCard(*state.Hand(), chosen)
	return &chosen
}

func FilterCards(state *State) []Card {
	ret := make([]Card, 0, len(*state.Hand()))
	for _, c := range *state.Hand() {
		if CheckCard(c, state) {
			if developing && debug {
				log.Printf("%v -> %#v OK", state.Print(), c.Print())
			}
			ret = append(ret, c)
		} else {
			if developing && debug {
				log.Printf("%v -> %#v NG", state.Print(), c.Print())
			}
		}
	}
	if developing && debug {
		log.Printf("ret: %v", PrintCards(ret))
		//log.Printf("%v", ret)
	}
	return ret
}

func CheckCard(chosen Card, state *State) bool {
	// 選んだカードのチェック
	if chosen.IsNum() {
		if chosen.Color() == state.Color {
			if developing && debug {
				log.Printf("c:%v %v, s.c:%v -> %v", chosen.Print(), chosen.Color(), state.Color, chosen.Color() == state.Color)
			}
			return true
		}
		if chosen.Num() == state.Number {
			if developing && debug {
				log.Printf("c:%v %v, s.c:%v -> %v", chosen.Print(), chosen.Color(), state.Color, chosen.Num() == state.Number)
			}
			return true
		}
		if developing && debug {
			log.Printf("c:%v %v, s.c:%v -> false", chosen.Print(), chosen.Color(), state.Color)
		}
	}
	if chosen.IsAction() {
		action := chosen.ActionType() + 100
		if developing && debug {
			log.Printf("chosen.Action: %v, action: %v", chosen.ActionType(), action)
		}
		switch action {
		case WILD___: // 場のカードがどのカードであっても出せる。
			if state.LastCard == nil { // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
				return true
			}
			if developing && debug {
				l := peek(&state.Discard)
				log.Printf("c:%v --, s:%v -> true", chosen.Print(), l.Print())
			}
			return true
		case DRAW4__:
			if state.LastCard == nil { // 場札 nilの時は直前の人がスキップしている。nil出ないときDiscardの末尾。
				return true //
			} else {
				// スキップされていないとき。
				// チャレンジルールはここではないところで判定する。
				// 前の人のカードがドローカードの時、ルール毎に出せるかどうか決まる。
				discard := peek(&state.Discard)
				action := discard.ActionType() + 100
				if action == DRAW2__ && action == DRAW2__ {
					if state.Rule.RequireMatchingDrawTwoForChallenge && action == DRAW2__ { // R.jua.05      ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
						if developing && debug {
							log.Printf("RequireMatchingDrawTwoForChallenge: %v : R.jua.05 ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。", state.Rule.RequireMatchingDrawTwoForChallenge)
						}
						return state.Color == discard.Color()
					}
					return state.Rule.StackDrawCards // R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
				}
				return true // ドローカード以外は場のカードがどのカードであっても出せる。
			}
		case DRAW2__:
			fallthrough
		case REVERSE:
			fallthrough
		case SKIP___:
			fallthrough
		case DISCARD:
			if state.Color == -1 {
				return true
			}
			return chosen.Color() == state.Color // 色が一致する必要がある。
		default:
			// 来ないはず。
			if developing && debug {
				log.Printf("chosen: %v", chosen)
			}
			panic(errors.Errorf("%v", "バグ"))
		}
	}
	return false // だめ。
}
