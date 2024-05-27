package core

import (
	"log"

	"github.com/pkg/errors"
)

/*

// カードを引く。
func DealCards(player_id int, deal_num int, state *State) bool {

	if state.Rule.EnableRandomDraw { // R.official-attack.02 引札の数をランダムにする。
		deal_num = rand.Intn(min(0, state.Rule.MaxDrawCount)) // R.official-attack.03 引札の最大枚数。5枚くらい?
	}
	if state.Rule.LimitHandSize { // R.jua.01 手札の枚数を制限する。
		hand_size := len(*state.Hand())
		max_hand_size := min(0, state.Rule.MaxHandSize) // R.jua.02 手札の枚数制限最大値。10枚。
		if deal_num+hand_size > max_hand_size {
			deal_num = max_hand_size - hand_size
		}
	}
	// deal_num回state.Deckからhandにカードを引く。
	for j := 0; j < deal_num; j++ {
		if len(state.Deck) == 0 {
			// 捨て札から場にシャッフルして戻す。
			last_card, ok := pop(&state.Discard) // 捨て札の最も新しいカードを保持する。
			if !ok {
				// 山札も捨て札もない。
				// 引けないので引かないというのもありかなと思われる。
				return false
			}
			ShuffleCards(state.Discard) // 捨て札をシャッフル。
			CryptoRandShuffle(state.Discard)
			state.Deck = append(state.Deck, state.Discard...) // 山札にする。
			state.Discard = []Card{last_card}                 // 捨て札に保持しておいた最後の捨て札を戻す。
		}
		// 山札からカードを1枚引く。
		if len(state.Deck) == 0 && len(state.Discard) <= 1 {
			// 山札も捨て札もない。
			// 引けないので引かないというのもありかなと思われる。
			if Developing && Debug {
				log.Printf("Deck: %v", PrintCards(state.Deck))
				log.Printf("Disc: %v", PrintCards(state.Discard))
				log.Printf("Hand: %v", PrintCards(*state.Hand()))
			}
			return false
		}
		card, ok := pop(&state.Deck)
		push(&state.Hands[player_id], card)
		if !ok {
			// 謎。バグ。
			panic(errors.Errorf("%v", "山札も捨て札もない。全てのカードが手札にある???"))
		}
	}
	if Developing && Debug {
		log.Printf("Dealed Cards : %#v", PrintCards(state.Hands[player_id]))
	}
	return true
}*/

func DealCards(player_id int, deal_num int, state *State) bool {
	state.DrawnCardID = []int{}
	// deal_num回state.Deckからhandにカードを引く。
	for j := 0; j < deal_num; j++ {
		// 山札からカードを1枚引く。
		if !CheckDeck(state) { // 山札が弾ける状態かチェック
			// 引けない状態なら引かない。
			continue
		}
		card, ok := pop(&state.Deck)
		if !ok {
			panic(errors.Errorf("%v", "バグ")) // 謎。バグ。
		}
		push(&state.Hands[player_id], card)
		state.DrawnCardID = append(state.DrawnCardID, card.Type)
	}
	if Developing && Debug {
		log.Printf("Dealed Cards : %#v", PrintCards(state.Hands[player_id]))
	}
	return true
}

func CheckDeck(state *State) bool {
	// 山札がない場合に、捨て札から場にシャッフルして戻す。
	ShuffleDiscardPileIntoDeck(state)

	if len(state.Deck) == 0 && len(state.Discard) <= 1 {
		// 山札も捨て札もない。
		// 引けないので引かないというのもありかなと思われる。
		if Developing && Debug {
			log.Printf("Deck: %v", PrintCards(state.Deck))
			log.Printf("Disc: %v", PrintCards(state.Discard))
			log.Printf("Hand: %v", PrintCards(*state.Hand()))
		}
		return false
	}
	return true
}

// 山札がない場合に、捨て札から場にシャッフルして戻す。
func ShuffleDiscardPileIntoDeck(state *State) {
	if len(state.Deck) == 0 && len(state.Discard) <= 1 {
		// 山札も捨て札もない。場合は何もしない。
		return
	}
	if len(state.Deck) == 0 && len(state.Discard) > 1 {
		// 捨て札から場にシャッフルして戻す。
		last_card, ok := pop(&state.Discard) // 捨て札の最も新しいカードを保持する。
		if !ok {
			// 来ないはず。
			panic(errors.Errorf("%v", "バグ"))
			return
		}
		ShuffleCards(state.Discard) // 捨て札をシャッフル。
		CryptoRandShuffle(state.Discard)
		state.Deck = append(state.Deck, state.Discard...) // 山札にする。
		state.Discard = []Card{last_card}                 // 捨て札に保持しておいた最後の捨て札を戻す。
	}
}
