package core

import "log"

func GetCards(rule UnoRule) []Card {
	ci := CardInfo{}
	ci.Init()

	if Developing && Debug {
		log.Printf("%v", rule.Print())
	}

	if rule.UseOfficialAppWildRules { // R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
		// 公式アプリ内のGOワイルドで使用されるカード
		// * 公式ワイルドルール
		// * 2倍のデッキ。
		// * +2に+2,+4を出せる。
		// * アクションカード枚数増加。ディスカードオールの追加。同じ色のカードをすべて捨てられる。
		//   アクションカード枚数増加と書いてあるがどれだけ増えているのか不明。
		//   元の112枚中36(8+4+8+8+8)枚の32%がアクションカードなので、4割程度になるように増加させると想定する。
		//   2倍のデッキの224(112x2)枚中アクションカードは72(36x2)枚で、4割は90枚程度なので20(90-72)枚弱増やす場合、
		//   元のアクションカードの半分の18(4+2+4+4+4)が妥当と思われる。
		//   また、これにディスカードオールをほかのアクションカード同じ12枚入れると想定すると、
		//   数字カード152(19x4x2)枚、アクションカード102(12+6+12+12+12+12)枚の合計254枚。
		//   これで、102/152=0.401なのでちょうど4割になる。

		// 2倍のデッキ
		for i, _ := range ci.CardTypeCount {
			ci.CardTypeCount[i] *= 2
		}
		// アクションカード増加
		// ディスカードオール増加 各色3枚 // R.official-app.08 ディスカードオールカードの追加枚数。各色枚数。
		ci.SetNumberOfCard("bDIS", 3)
		ci.SetNumberOfCard("gDIS", 3)
		ci.SetNumberOfCard("rDIS", 3)
		ci.SetNumberOfCard("yDIS", 3)

		ci.AddNumberOfCard("WILD", 4) // ワイルド // 4枚
		ci.AddNumberOfCard("WD+4", 2) // ドロー4 // 2枚
		ci.AddNumberOfCard("bD+2", 1) // ドロー2 // 各色1枚
		ci.AddNumberOfCard("gD+2", 1)
		ci.AddNumberOfCard("rD+2", 1)
		ci.AddNumberOfCard("yD+2", 1)
		ci.AddNumberOfCard("bREV", 1) // リバース // 各色1枚
		ci.AddNumberOfCard("gREV", 1)
		ci.AddNumberOfCard("rREV", 1)
		ci.AddNumberOfCard("yREV", 1)
		ci.AddNumberOfCard("bSKP", 1) // スキップ // 各色1枚
		ci.AddNumberOfCard("gSKP", 1)
		ci.AddNumberOfCard("rSKP", 1)
		ci.AddNumberOfCard("ySKP", 1)
	}
	ci.SetNumberOfCard("bDIS", rule.DiscardAllExtraCount) // ディスカードオールの追加。
	ci.SetNumberOfCard("gDIS", rule.DiscardAllExtraCount)
	ci.SetNumberOfCard("rDIS", rule.DiscardAllExtraCount)
	ci.SetNumberOfCard("yDIS", rule.DiscardAllExtraCount)

	ci.SetNumberOfCard("Shff", rule.ShuffleWildExtraCount) // シャッフルワイルドカードの追加 シャッフルワイルド有りの時、ワイルドカードの枚数が変わる。

	// R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
	// ホワイトワイルドをどうするか。とりあえずワイルドカードにしておく。 追々実装してもよい。
	if rule.WhiteWildCount != 0 {
		ci.SetNumberOfCard("WILD", 4+rule.WhiteWildCount) // ワイルド // 8枚 -> 7枚 (ワイルド4枚 ホワイトワイルド3枚)
	}

	{
		// R.official-attack.01 公式アプリのアタックエクストリームルールを使用する。
		if rule.UseOfficialAttackExtreamRules {
			// 公式アプリのアタックエクストリームルールを使用する。
			// 既存カードの枚数が変わる。
			ci.SetNumberOfCard("WILD", 4) // ワイルド         4枚 (-4)
			ci.AddNumberOfCard("WILD", 3) // ホワイトワイルド 3枚
			ci.AddNumberOfCard("WD+4", 0) // ワイルドドロー4  0枚 (-4)
			ci.SetNumberOfCard("bREV", 1) // リバース 各色1枚 4枚 (-4)
			ci.SetNumberOfCard("gREV", 1) //
			ci.SetNumberOfCard("rREV", 1) //
			ci.SetNumberOfCard("yREV", 1) //
			ci.SetNumberOfCard("bSKP", 2) // スキップ 各色2枚 8枚
			ci.SetNumberOfCard("gSKP", 2) //
			ci.SetNumberOfCard("rSKP", 2) //
			ci.SetNumberOfCard("ySKP", 2) //
			ci.SetNumberOfCard("bDIS", 2) // ディスカードオール 各色2枚 8枚 (+8)
			ci.SetNumberOfCard("gDIS", 2) //
			ci.SetNumberOfCard("rDIS", 2) //
			ci.SetNumberOfCard("yDIS", 2) //
			ci.SetNumberOfCard("bHT2", 2) // ヒット2 各色2枚  8枚 (+8)
			ci.SetNumberOfCard("gHT2", 2) //
			ci.SetNumberOfCard("rHT2", 2) //
			ci.SetNumberOfCard("yHT2", 2) //
			ci.SetNumberOfCard("WHT4", 1) // ワイルドヒット4                1枚 (+1)
			ci.SetNumberOfCard("WAA2", 4) // ワイルドアタックアタックカード 4枚 (+4)
		}
		ci.AddNumberOfCard("WILD", rule.WhiteWildCount)     // ホワイトワイルド
		ci.AddNumberOfCard("bHT2", rule.Hit2ExtraCount)     // ヒット2カードの追加
		ci.AddNumberOfCard("gHT2", rule.Hit2ExtraCount)     //
		ci.AddNumberOfCard("rHT2", rule.Hit2ExtraCount)     //
		ci.AddNumberOfCard("yHT2", rule.Hit2ExtraCount)     //
		ci.AddNumberOfCard("WHT4", rule.WildHit4ExtraCount) // ワイルド ヒット4の追加
		ci.AddNumberOfCard("WAA2", rule.WildHit4ExtraCount) // ワイルド アタックアタックカードの追加
	}

	// 日本UNO協会ルール。
	if rule.FlatRateNumberCardPoints { // 数字カードのポイント計算を一律に5点としてカウントする。
		for _, name := range ci.CardNameList {
			c := ci.Get(name)
			if c.IsNum() {
				c.Point = 5
			}
		}
	}

	return ci.GenerateDeck()
}
