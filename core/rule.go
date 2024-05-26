package core

import "fmt"

var (
	UnoRuleList = []string{
		"official",     //: 旧公式ルール。112枚。シャッフルワイルドなし。(ホワイトワイルドは未実装。)
		"official_new", //: 新公式ルール。シャッフルワイルドあり。(ホワイトワイルドは未実装。)
		//:   - シャッフルワイルドカードの追加枚数
		"official_app_classic", //: 公式アプリルール。公式ルールに時間制限と自動チャレンジ判定が追加。
		//:   - ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
		//:   - 1ターンの制限時間(秒)。0の時制限しない。
		//:   - 1ゲームの制限時間(秒)。0の時制限しない。
		"official_app_wild", //: 公式アプリワイルドルール。公式アプリルールベース。
		//:   - ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
		//:   - ディスカードオールカードの追加枚数。各色枚数。
		"official_app_classic_auto_uno", //: 公式アプリ課金ルール。公式アプリルールかつUNOコールが不要になる。
		"official_app_wild_auto_uno",    //: 公式アプリワイルド課金ルール。公式アプリワイルドルールかつ課金でUNOコールが不要になる。
		"official_attack",               //: 公式アタックエクストリームルール。
		//:   - 引札の数をランダムにする。
		//:   - 引札の最大枚数。5枚くらい?
		//:   - ヒット2カードの追加枚数。
		//:   - ワイルドヒット4の追加枚数。
		//:   - ワイルドアタックアタックカードの追加枚数。
		"jp_uno_offical", //: 日本UNO協会が定めた競技ルール。公式ルールをベースにする。
		//:   - 手札の枚数を制限する。
		//:   - 手札の枚数制限最大値。10枚。
		//:   - 同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		//:   - 数字カードのポイント計算を一律に5点としてカウントする。
		//:   - ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
		"algori", //: ALGORI大会ルール。新公式ルールをベースにする。
		//:   - ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
		//:   - 手札の枚数を制限する。
		//:   - 手札の枚数制限最大値。25枚。
		//:   - 1ターンの制限時間(5秒)。0の時制限しない。
		//:   - 不正操作時にカードを戻し、2枚引き、ターンをスキップするルールを有効にする。
		//:   - 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
		//:   - ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す。
		"house-jp", //: 日本の津々浦々で遊ばれているルール。色々ある。新公式ルールをベースにする。
		//:   - 記号カードでの勝利を禁止する。
		//:   - ゲームを最後の1人まで継続する。
		//:   - 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
		//:   - 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		//:   - 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		//:   - 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
		"custom1", //: ユーザーが好きにカスタマイズできるルール。
		"custom2", //: ユーザーが好きにカスタマイズできるルール。追加してもよい。
	} // 実装予定ルールリスト
)

func (rule *UnoRule) Init() UnoRule {
	// 公式ルール。
	rule.NumberOfPlayer = 4 // R.official.01 人数は2-10人程度を許容。4-6人を最適とする。

	// 公式ルールではないが、このプログラムで規定とするルール。
	rule.EveryoneLosesAfterConsecutiveSkips = 10 // 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
	rule.PlayFromHandAfterDraw = false           // (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。

	// 書いていないが公式ルール
	// ドロー4にドロー4が重ねられたときチャレンジを禁止する。そもそも重ねられない。重ねられる場合ドロー2→ドロー4の時のみのルール。このルールはカスタマイズできなくてよいと思われる。

	return *rule
}

func (ur *UnoRule) SetRule(rule string) {

	switch rule {
	case "official": // 公式ルール。
		*ur = (&UnoRule{}).Init() // 全て初期化する。
	case "official_new": // シャッフルワイルドあり。
		ur.SetRule("official")       // 公式ルールをベースにする。
		ur.ShuffleWildExtraCount = 1 // R.official.02 シャッフルワイルドカードの追加枚数
		ur.WhiteWildCount = 3        // R.official.03 ホワイトワイルドカードの追加枚数
	case "official_app_classic": // 公式アプリ基本ルール。
		ur.SetRule("official")               // 公式ルールをベースにする。
		ur.AutomatedChallengeJudgment = true // R.official-app.01 ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
		ur.TurnTimeLimit = 10                // R.official-app.02 1ターンの制限時間(秒)。0の時制限しない。
		ur.GameTimeLimit = 180               // R.official-app.03 1ゲームの制限時間(秒)。0の時制限しない。
	case "official_app_wild": // 公式アプリワイルドルール。
		ur.SetRule("official_app_classic") // 公式アプリクラシックルールをベースにする。
		ur.UseOfficialAppWildRules = true  // R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
		ur.StackDrawCards = true           // R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
		//ur.DiscardAllExtraCount = 0        // R.official-app.08 ディスカードオールカードの追加枚数。各色枚数。

	case "official_app_classic_auto_uno": // 公式基本ルール。課金でUNOコールが不要になる。
		ur.SetRule("official_app_classic") // 公式アプリクラシックルールをベースにする。
		ur.AutoUnoCall = true              // R.official-app.04 UNOコールを不要とする。 「ウノ」と叫ぶのを忘れた人が最後から2枚目のカードを捨てた瞬間から、次の人がカードを捨てる瞬間までの間に指摘されたら、罰になります。
	case "official_app_wild_auto_uno": // 公式アプリワイルドルール。課金でUNOコールが不要になる。
		ur.SetRule("official_app_wild") // 公式アプリクラシックルールをベースにする。
		ur.AutoUnoCall = true           // R.official-app.04 UNOコールを不要とする。 「ウノ」と叫ぶのを忘れた人が最後から2枚目のカードを捨てた瞬間から、次の人がカードを捨てる瞬間までの間に指摘されたら、罰になります。
	case "official_attack": // 公式アタックエクストリームルール。罰則で引く枚数がランダム。https://mattel.co.jp/toys/mattel_games/mattel_games-11044/#howToPlay
		ur.UseOfficialAttackExtreamRules = true // R.official-attack.01 公式アプリのアタックエクストリームルールを使用する。
		ur.EnableRandomDraw = true              // R.official-attack.02 引札の数をランダムにする。
		ur.MaxDrawCount = 5                     // R.official-attack.03 引札の最大枚数。5枚くらい?
		ur.Hit2ExtraCount = 4                   // R.official-attack.04 ヒット2カードの追加枚数。
		ur.WildHit4ExtraCount = 1               // R.official-attack.05 ワイルドヒット4の追加枚数。
		ur.WildAttackAttackExtraCount = 4       // R.official-attack.06 ワイルドアタックアタックカードの追加枚数。
	case "jp_uno_offical": // 日本UNO協会が定めた競技ルール。古いリンクは死んでいる。https://web.archive.org/web/20190331140300/http://www.geocities.jp/unoassoc/rule.html https://github.com/m0bec/UNO/issues/6
		ur.SetRule("official")                       // 公式ルールをベースにする。
		ur.LimitHandSize = true                      // R.jua.01      手札の枚数を制限する。
		ur.MaxHandSize = 10                          // R.jua.02      手札の枚数制限最大値。10枚。
		ur.AllowSimultaneousSameColor = true         // R.jua.03      同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.FlatRateNumberCardPoints = true           // R.jua.04      数字カードのポイント計算を一律に5点としてカウントする。
		ur.RequireMatchingDrawTwoForChallenge = true // R.jua.05      ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
	case "algori":
		// ALGORI 大会ルール https://www.d3.ntt-east.co.jp/algori/rule.html
		// 上のサイトを読む時以下を前提にしておかないと読みにくい。文言が曖昧。独自用語あり。実質的にどちらでもよいルールは書いていない感じ。
		// 書いてない部分は公式ルール準拠で、書いてあるルールはすべて公式ルールと異なる部分だと思って読む必要がある。
		// - プレーヤー引くのではなく、ディーラーが配る。ルールを読む時の視点が異なる。
		// - ディーラー(親ではない。おそらく親がいない。)がプレーヤーとは別にいる。ディーラー1人、プレーヤー4人の5人いると見做せる。ディーラーはプログラム。
		// - ゲームの開始時点のみ、5人目のディーラーがプレーヤーに加わっている状態を想定して読むとわかりやすいかもしれない。
		// 噛み砕いたルール
		// - 順番(回る方向)と最初のプレーヤーはランダム。
		// - 最初のカードの決定方法が特殊。
		//   - ディーラーが1枚出す。場札(最初の場のカードを決めるためにディーラーが配ったカード)とする。
		//   - 場札が
		//     - 色が決まらないカード(ワイルドドロー4、シャッフルワイルド、ホワイトワイルド)の場合は、山札のランダムな位置に戻して、再度ディーラーが場札を出す。
		//     - スキップの場合、スキップの効果が適用される。2番目のプレーヤーからスタート。
		//     - ドロー2の時、ドロー2の効果が適用される。1番目のプレーヤーが2枚引き、2番目のプレイヤーからスタート。
		//     - リバースの時、リバースの効果が適用される。3番目のプレーヤーからリバースしてスタート。
		//   - 場札が決まった後に手札を配る。(公式ルールと異なる。)
		// - 1ターンは5秒以内。
		// - 出せるカードの有無に関わらず(手札が25枚未満の時)カードを1枚引くことができるが、引いたときは引いたカードのみ出せる。
		// - 不正操作の場合(時間切れで出した、誤ったカードを出した、色を指定すべき時に指定しなかった、色を指定すべきでないとき指定した、手札が複数枚の時にUNOコールなど)
		//   - 罰則で2枚引く。出したカードは手札に戻る。スキップされる。
		// - 場札と山札両方が無くなった場合(全て手札にある状態。場札は1枚だと思われる。)、10周経過して誰も出さない場合、勝者なしとして全員敗者扱いする。
		// - ワイルドチャレンジが成功した時、ワイルドドロー4は手札に戻される。
		ur.SetRule("official_app_classic")                // 公式追加ルール
		ur.LimitHandSize = true                           // 手札の枚数を制限する。
		ur.MaxHandSize = 25                               // 手札の枚数制限最大値。25枚。
		ur.TurnTimeLimit = 5                              // R.official-app.03 1ターンの制限時間(秒)。0の時制限しない。
		ur.EnableAlgoriStart = true                       // R.algori.O1   ALGORIルールの開始方法を使用する。
		ur.EnforceCheatingPenalty = true                  // R.algori.O2   不正操作時にカードを戻し、2枚引き、ターンをスキップするルールが有効か
		ur.SkipRoundsGameOver = true                      // R.algori.O3   10周スキップした場合全員敗者とする
		ur.ReturnWildDrawFourOnSuccessfulChallenge = true // R.algori.O4   ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す。
		ur.EveryoneLosesAfterConsecutiveSkips = 10        // R.algori.O5   一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
	case "house-jp": // 日本の津々浦々で遊ばれているルール。色々ある。
		ur.PlayFromHandAfterDraw = true            // R.house-jp.01 (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。
		ur.RestrictWinOnSpecialCard = true         // R.house-jp.02 記号カードでの勝利を禁止する。
		ur.ContinueUntilLastPlayer = true          // R.house-jp.03 ゲームを最後の1人まで継続する。
		ur.DisallowImmediatePlayOfDrawCards = true // R.house-jp.04 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
		ur.AllowSequentialSameColor = true         // R.house-jp.05 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.AllowSimultaneousAnyColor = true        // R.house-jp.06 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.AllowWinWithoutUnoOnSimultaneous = true // R.house-jp.07 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
	case "custom1": // ユーザーが好きにカスタマイズできるルール。
		ur.UserCustumRule(1)
	case "custom2": // ユーザーが好きにカスタマイズできるルール。追加してもよい。
		ur.UserCustumRule(2)
	}
}

func (ur *UnoRule) UserCustumRule(num int) {
	// 設定ファイルを読み込む、など
	// 未実装
}

// カスタムできるようにルールをON/OFFできるようにする。
type UnoRule struct {

	// 公式ルールここから。 {{{

	// falseの状態が公式ルールになるようにする。
	NumberOfPlayer int `json:"number_of_player"` // R.official.01 人数は2-10人程度を許容。4-6人を最適とする。
	// 新公式ルール
	// 公式ルールも年毎にカードの追加などがある。https://mattel.co.jp/toys/mattel_games/mattel_games-10936/#howToPlay
	ShuffleWildExtraCount int `json:"shuffle_wild_extra_count"` // R.official.02 シャッフルワイルドカードの追加枚数
	WhiteWildCount        int `json:"white_wild_extra_count"`   // R.official.03 ホワイトワイルドカードの追加枚数
	// シャッフルワイルドカードを出す人は、全員のカードを集めてシャッフルしてください。
	// 自分の左隣の人から順番に１枚ずつ、すべて配ります。カードが増える人もいるし、減る人もいます。
	// 好きな色を宣言し、次の人に順番が移ります。場のカードが何であっても、捨てることができます。
	// 手持ちのカードの中に、他に使えるカードがあってもこのカードを使えます。
	// 最初の場のカードがこのカードだった時は、親の左どなりの人（最初のプレイヤー）が好きな色を宣言してカードを捨てます。

	// 公式アプリルール。
	AutomatedChallengeJudgment bool `json:"automated_challenge_judgment"` // R.official-app.01 ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
	TurnTimeLimit              int  `json:"turn_time_limit"`              // R.official-app.02 1ターンの制限時間(秒)。0の時制限しない。
	GameTimeLimit              int  `json:"game_time_limit"`              // R.official-app.03 1ゲームの制限時間(秒)。0の時制限しない。
	AutoUnoCall                bool `json:"auto_uno_call"`                // R.official-app.04 UNOコールを不要とする。 「ウノ」と叫ぶのを忘れた人が最後から2枚目のカードを捨てた瞬間から、次の人がカードを捨てる瞬間までの間に指摘されたら、罰になります。
	TwoVsTwoMode               bool `json:"two_vs_two_mode"`              // R.official-app.05 2対2の対戦モードを有効にする

	// 出せるカードの有無に関わらず(手札が25枚未満の時)カードを1枚引くことができるが、引いたときは引いたカードのみ出せる。
	// 公式アプリワイルドルール。
	UseOfficialAppWildRules bool `json:"use_official_app_wild_rules"` // R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
	StackDrawCards          bool `json:"stack_draw_cards"`            // R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
	DiscardAllExtraCount    int  `json:"discard_all_extra_count"`     // R.official-app.08 ディスカードオールカードの追加枚数。各色枚数。
	// ディスカードオールカードの追加枚数。
	// シャッフルワイルドカードの追加枚数。
	// ヒット2カードの追加枚数。
	// ワイルドヒット4の追加枚数。
	// ワイルドアタックアタックカードの追加枚数。

	// ウノ アタック エクストリームルール。 https://mattel.co.jp/toys/mattel_games/mattel_games-11044/#howToPlay
	UseOfficialAttackExtreamRules bool `json:"use_official_attack_extream_rules"` // R.official-attack.01 公式アプリのアタックエクストリームルールを使用する。
	EnableRandomDraw              bool `json:"enable_random_draw"`                // R.official-attack.02 引札の数をランダムにする。
	MaxDrawCount                  int  `json:"max_draw_count"`                    // R.official-attack.03 引札の最大枚数。5枚くらい?
	Hit2ExtraCount                int  `json:"hit_two_extra_count"`               // R.official-attack.04 ヒット2カードの追加枚数。
	WildHit4ExtraCount            int  `json:"wild_hit_four_extra_count"`         // R.official-attack.05 ワイルドヒット4の追加枚数。
	WildAttackAttackExtraCount    int  `json:"wild_attack_attack_extra_count"`    // R.official-attack.06 ワイルドアタックアタックカードの追加枚数。
	// 場のカードが何であっても、捨てることができます。
	// 好きな色を宣言し、アタック攻撃をしかける相手を決めます。
	// アタック攻撃を受けた人はアタックボタンを2回押します。
	// 注意:ワイルド アタックアタックを出してあがる時は、通常どおり好きな色を宣言し、アタック攻撃をしかける相手を決めてあがってください。(点数計算のため)
	// 最初の場のカードがこのカードだった時は、親の左どなりの人(最初のプレイヤー)がアタック攻撃をしかける相手を決めます。
	// アタック攻撃を受けた人はアタックボタンを2回押します。
	// アタック攻撃を受けた人の次の人(親の左側の二人目のプレイヤー) に順番が移ります。
	// 親の左どなりの人最初のプレイヤー)が好きな色を宣言します。

	// 公式ルールここまで。 }}}

	// 日本UNO協会ルール。
	// 古い一覧は死んでいる。https://web.archive.org/web/20190331140300/http://www.geocities.jp/unoassoc/rule.html
	LimitHandSize                      bool `json:"limit_hand_size"`                         // R.jua.01 手札の枚数を制限する。
	MaxHandSize                        int  `json:"max_hand_size"`                           // R.jua.02 手札の枚数制限最大値。10枚。
	AllowSimultaneousSameColor         bool `json:"allow_simultaneous_same_color"`           // R.jua.03 同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	FlatRateNumberCardPoints           bool `json:"flat_rate_number_card_points"`            // R.jua.04 数字カードのポイント計算を一律に5点としてカウントする。
	RequireMatchingDrawTwoForChallenge bool `json:"require_matching_draw_two_for_challenge"` // R.jua.05 ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。

	// ALGORI 大会ルール https://www.d3.ntt-east.co.jp/algori/rule.html
	// 上のサイトを読む時以下を前提にしておかないと読みにくい。文言が曖昧。独自用語あり。実質的にどちらでもよいルールは書いていない感じ。
	// 書いてない部分は公式ルール準拠で、書いてあるルールはすべて公式ルールと異なる部分だと思って読む必要がある。
	// - プレーヤー引くのではなく、ディーラーが配る。ルールを読む時の視点が異なる。
	// - ディーラー(親ではない。おそらく親がいない。)がプレーヤーとは別にいる。ディーラー1人、プレーヤー4人の5人いると見做せる。ディーラーはプログラム。
	// - ゲームの開始時点のみ、5人目のディーラーがプレーヤーに加わっている状態を想定して読むとわかりやすいかもしれない。
	// 噛み砕いたルール
	// - 順番(回る方向)と最初のプレーヤーはランダム。
	// - 最初のカードの決定方法が特殊。
	//   - ディーラーが1枚出す。場札(最初の場のカードを決めるためにディーラーが配ったカード)とする。
	//   - 場札が
	//     - 色が決まらないカード(ワイルドドロー4、シャッフルワイルド、ホワイトワイルド)の場合は、山札のランダムな位置に戻して、再度ディーラーが場札を出す。
	//     - スキップの場合、スキップの効果が適用される。2番目のプレーヤーからスタート。
	//     - ドロー2の時、ドロー2の効果が適用される。1番目のプレーヤーが2枚引き、2番目のプレイヤーからスタート。
	//     - リバースの時、リバースの効果が適用される。3番目のプレーヤーからリバースしてスタート。
	//   - 場札が決まった後に手札を配る。(公式ルールと異なる。)
	// - 1ターンは5秒以内。
	// - 出せるカードの有無に関わらず(手札が25枚未満の時)カードを1枚引くことができるが、引いたときは引いたカードのみ出せる。
	// - 不正操作の場合(時間切れで出した、誤ったカードを出した、色を指定すべき時に指定しなかった、色を指定すべきでないとき指定した、手札が複数枚の時にUNOコールなど)
	//   - 罰則で2枚引く。出したカードは手札に戻る。スキップされる。
	// - 場札と山札両方が無くなった場合(全て手札にある状態。場札は1枚だと思われる。)、10周経過して誰も出さない場合、勝者なしとして全員敗者扱いする。
	// - ワイルドチャレンジが成功した時、ワイルドドロー4は手札に戻される。

	EnableAlgoriStart                       bool `json:"enable_algori_start"`                    // R.algori.O1 ALGORIルールの開始方法を使用する。
	EnforceCheatingPenalty                  bool `json:"enforce_cheating_penalty"`               // R.algori.O2 不正操作時にカードを戻し、2枚引き、ターンをスキップするルールが有効か
	SkipRoundsGameOver                      bool `json:"skip_rounds_game_over"`                  // R.algori.O3 10周スキップした場合全員敗者とする
	ReturnWildDrawFourOnSuccessfulChallenge bool `json:"return_wild_draw_four_on_challenge"`     // R.algori.O4 ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す
	EveryoneLosesAfterConsecutiveSkips      int  `json:"everyone_loses_after_consecutive_skips"` // R.algori.O5 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。

	// 非公式ハウスルール。
	// 上記になく、日本のローカルルール集にあったもの。

	PlayFromHandAfterDraw            bool `json:"play_from_hand_after_draw"`             // R.house-jp.01 (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。
	RestrictWinOnSpecialCard         bool `json:"restrict_win_on_special_card"`          // R.house-jp.02 記号カードでの勝利を禁止する。
	ContinueUntilLastPlayer          bool `json:"continue_until_last_player"`            // R.house-jp.03 ゲームを最後の1人まで継続する。
	DisallowImmediatePlayOfDrawCards bool `json:"disallow_immediate_play_of_draw_cards"` // R.house-jp.04 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
	AllowSequentialSameColor         bool `json:"allow_sequential_same_color"`           // R.house-jp.05 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	AllowSimultaneousAnyColor        bool `json:"allow_simultaneous_any_color"`          // R.house-jp.06 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	AllowWinWithoutUnoOnSimultaneous bool `json:"allow_win_without_uno_on_simultaneous"` // R.house-jp.07 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
}

func (ur *UnoRule) Print() string {
	return fmt.Sprintf(`
	--------------------------------------------------------------------------------------------------------------------------
	NumberOfPlayer,                         | %5v | R.official.01 人数は2-10人程度を許容。4-6人を最適とする。
	ShuffleWildExtraCount,                  | %5v | R.official.02 シャッフルワイルドカードの追加枚数
	WhiteWildCount,                         | %5v | R.official.03 ホワイトワイルドカードの追加枚数
	--------------------------------------------------------------------------------------------------------------------------
	AutomatedChallengeJudgment,             | %5v | R.official-app.01 ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
	TurnTimeLimit,                          | %5v | R.official-app.02 1ターンの制限時間(秒)。0の時制限しない。
	GameTimeLimit,                          | %5v | R.official-app.03 1ゲームの制限時間(秒)。0の時制限しない。
	AutoUnoCall,                            | %5v | R.official-app.04 UNOコールを不要とする。 「ウノ」と叫ぶのを忘れた人が最後から2枚目のカードを捨てた瞬間から、次の人がカードを捨てる瞬間までの間に指摘されたら、罰になります。
	TwoVsTwoMode,                           | %5v | R.official-app.05 2対2の対戦モードを有効にする
	UseOfficialAppWildRules,                | %5v | R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
	StackDrawCards,                         | %5v | R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
	DiscardAllExtraCount,                   | %5v | R.official-app.08 ディスカードオールカードの追加枚数。各色枚数。
	--------------------------------------------------------------------------------------------------------------------------
	UseOfficialAttackExtreamRules,          | %5v | R.official-attack.01 公式アプリのアタックエクストリームルールを使用する。
	EnableRandomDraw,                       | %5v | R.official-attack.02 引札の数をランダムにする。
	MaxDrawCount,                           | %5v | R.official-attack.03 引札の最大枚数。5枚くらい?
	Hit2ExtraCount,                         | %5v | R.official-attack.04 ヒット2カードの追加枚数。
	WildHit4ExtraCount,                     | %5v | R.official-attack.05 ワイルドヒット4の追加枚数。
	WildAttackAttackExtraCount,             | %5v | R.official-attack.06 ワイルドアタックアタックカードの追加枚数。
	----------------------------------------------------------------------------------------------------------------------------
	LimitHandSize                           | %5v | R.jua.01      | 手札の枚数を制限する。
	MaxHandSize                             | %5v | R.jua.02      | 手札の枚数制限最大値。10枚。
	AllowSimultaneousSameColor              | %5v | R.jua.03      | 同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	FlatRateNumberCardPoints                | %5v | R.jua.04      | 数字カードのポイント計算を一律に5点としてカウントする。
	RequireMatchingDrawTwoForChallenge      | %5v | R.jua.05      | ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
	----------------------------------------------------------------------------------------------------------------------------
	EnableAlgoriStart                       | %5v | R.algori.O1   | ALGORIルールの開始方法を使用する。
	EnforceCheatingPenalty                  | %5v | R.algori.O2   | 不正操作時にカードを戻し、2枚引き、ターンをスキップするルールが有効か
	SkipRoundsGameOver                      | %5v | R.algori.O3   | 10周スキップした場合全員敗者とする
	ReturnWildDrawFourOnSuccessfulChallenge | %5v | R.algori.O4   | ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す
	EveryoneLosesAfterConsecutiveSkips      | %5v | R.algori.O5   | 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
	----------------------------------------------------------------------------------------------------------------------------
	PlayFromHandAfterDraw                   | %5v | R.house-jp.01 | (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。
	RestrictWinOnSpecialCard                | %5v | R.house-jp.02 | 記号カードでの勝利を禁止する。
	ContinueUntilLastPlayer                 | %5v | R.house-jp.03 | ゲームを最後の1人まで継続する。
	DisallowImmediatePlayOfDrawCards        | %5v | R.house-jp.04 | 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
	AllowSequentialSameColor                | %5v | R.house-jp.05 | 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	AllowSimultaneousAnyColor               | %5v | R.house-jp.06 | 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	AllowWinWithoutUnoOnSimultaneous        | %5v | R.house-jp.07 | 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
	----------------------------------------------------------------------------------------------------------------------------
	`,
		ur.NumberOfPlayer,                          // R.official.01 人数は2-10人程度を許容。4-6人を最適とする。
		ur.ShuffleWildExtraCount,                   // R.official.02 シャッフルワイルドカードの追加枚数
		ur.WhiteWildCount,                          // R.official.03 ホワイトワイルドカードの追加枚数
		ur.AutomatedChallengeJudgment,              // R.official-app.01 ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
		ur.TurnTimeLimit,                           // R.official-app.02 1ターンの制限時間(秒)。0の時制限しない。
		ur.GameTimeLimit,                           // R.official-app.03 1ゲームの制限時間(秒)。0の時制限しない。
		ur.AutoUnoCall,                             // R.official-app.04 UNOコールを不要とする。 「ウノ」と叫ぶのを忘れた人が最後から2枚目のカードを捨てた瞬間から、次の人がカードを捨てる瞬間までの間に指摘されたら、罰になります。
		ur.TwoVsTwoMode,                            // R.official-app.05 2対2の対戦モードを有効にする
		ur.UseOfficialAppWildRules,                 // R.official-app.06 公式アプリのワイルドルールを使用する。カード枚数が増える。
		ur.StackDrawCards,                          // R.official-app.07 ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
		ur.DiscardAllExtraCount,                    // R.official-app.08 ディスカードオールカードの追加枚数。各色枚数。
		ur.UseOfficialAttackExtreamRules,           // R.official-attack.01 公式アプリのアタックエクストリームルールを使用する。
		ur.EnableRandomDraw,                        // R.official-attack.02 引札の数をランダムにする。
		ur.MaxDrawCount,                            // R.official-attack.03 引札の最大枚数。5枚くらい?
		ur.Hit2ExtraCount,                          // R.official-attack.04 ヒット2カードの追加枚数。
		ur.WildHit4ExtraCount,                      // R.official-attack.05 ワイルドヒット4の追加枚数。
		ur.WildAttackAttackExtraCount,              // R.official-attack.06 ワイルドアタックアタックカードの追加枚数。
		ur.LimitHandSize,                           // R.jua.01      手札の枚数を制限する。
		ur.MaxHandSize,                             // R.jua.02      手札の枚数制限最大値。10枚。
		ur.AllowSimultaneousSameColor,              // R.jua.03      同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.FlatRateNumberCardPoints,                // R.jua.04      数字カードのポイント計算を一律に5点としてカウントする。
		ur.RequireMatchingDrawTwoForChallenge,      // R.jua.05      ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
		ur.EnableAlgoriStart,                       // R.algori.O1   ALGORIルールの開始方法を使用する。
		ur.EnforceCheatingPenalty,                  // R.algori.O2   不正操作時にカードを戻し、2枚引き、ターンをスキップするルールが有効か
		ur.SkipRoundsGameOver,                      // R.algori.O3   10周スキップした場合全員敗者とする
		ur.ReturnWildDrawFourOnSuccessfulChallenge, // R.algori.O4   ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す
		ur.EveryoneLosesAfterConsecutiveSkips,      // R.algori.O5   一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
		ur.PlayFromHandAfterDraw,                   // R.house-jp.01 (ドローの効果以外で)引いた場合でも元の手札から出すことを許可。
		ur.RestrictWinOnSpecialCard,                // R.house-jp.02 記号カードでの勝利を禁止する。
		ur.ContinueUntilLastPlayer,                 // R.house-jp.03 ゲームを最後の1人まで継続する。
		ur.DisallowImmediatePlayOfDrawCards,        // R.house-jp.04 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
		ur.AllowSequentialSameColor,                // R.house-jp.05 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.AllowSimultaneousAnyColor,               // R.house-jp.06 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
		ur.AllowWinWithoutUnoOnSimultaneous,        // R.house-jp.07 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
	)
}
