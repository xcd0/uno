package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/hjson/hjson-go/v4"
	"github.com/pkg/errors"
)

type Setting struct {
	// RuleToApply に指定できる文字列
	// "official"                      : 旧公式ルール。112枚。シャッフルワイルドなし。(ホワイトワイルドは未実装。)
	// "official_new"                  : 新公式ルール。シャッフルワイルドあり。(ホワイトワイルドは未実装。)
	//                                 :   - シャッフルワイルドカードの追加枚数
	// "official_app_classic"          : 公式アプリルール。公式ルールに時間制限と自動チャレンジ判定が追加。
	//                                 :   - ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
	//                                 :   - 1ターンの制限時間(秒)。0の時制限しない。
	//                                 :   - 1ゲームの制限時間(秒)。0の時制限しない。
	// "official_app_wild"             : 公式アプリワイルドルール。公式アプリルールベース。
	//                                 :   - ドロー2にドロー2やドロー4、ドロー4にドロー4を重ねることができる。
	//                                 :   - ディスカードオールカードの追加枚数。各色枚数。
	// "official_app_classic_auto_uno" : 公式アプリ課金ルール。公式アプリルールかつUNOコールが不要になる。
	// "official_app_wild_auto_uno"    : 公式アプリワイルド課金ルール。公式アプリワイルドルールかつ課金でUNOコールが不要になる。
	// "official_attack"               : 公式アタックエクストリームルール。
	//                                 :   - 引札の数をランダムにする。
	//                                 :   - 引札の最大枚数。5枚くらい?
	//                                 :   - ヒット2カードの追加枚数。
	//                                 :   - ワイルドヒット4の追加枚数。
	//                                 :   - ワイルドアタックアタックカードの追加枚数。
	// "jp_uno_offical"                : 日本UNO協会が定めた競技ルール。公式ルールをベースにする。
	//                                 :   - 手札の枚数を制限する。
	//                                 :   - 手札の枚数制限最大値。10枚。
	//                                 :   - 同じ数字かつ同じ色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	//                                 :   - 数字カードのポイント計算を一律に5点としてカウントする。
	//                                 :   - ドロー2にドロー4が重ねられたとき、チャレンジするためには手札に場の色と同じ色のドロー2を必要とする。
	// "algori"                        : ALGORI大会ルール。新公式ルールをベースにする。
	//                                 :   - ワイルドドロー4のチャレンジルールの際に、手札を見せて判定するのではなく、システムが判定する。
	//                                 :   - 手札の枚数を制限する。
	//                                 :   - 手札の枚数制限最大値。25枚。
	//                                 :   - 1ターンの制限時間(5秒)。0の時制限しない。
	//                                 :   - 不正操作時にカードを戻し、2枚引き、ターンをスキップするルールを有効にする。
	//                                 :   - 一定周スキップした場合全員敗者とするスキップ周回数。0の時無効。
	//                                 :   - ワイルドチャレンジが成功した時、ワイルドドロー4を手札に戻す。
	// "house-jp"                      : 日本の津々浦々で遊ばれているルール。色々ある。新公式ルールをベースにする。
	//                                 :   - 記号カードでの勝利を禁止する。
	//                                 :   - ゲームを最後の1人まで継続する。
	//                                 :   - 引いた直後のドロー２またはドロー４をすぐに出すことを禁止する。
	//                                 :   - 同じ色の数字の連番のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	//                                 :   - 同じ数字で任意の色のカードの同時出しルールを許可する。(上がる際にはUNOコールが必要、コールなしの場合1枚しか出せない、アクションカードは禁止)
	//                                 :   - 同時出しルールが許可されている状態で、UNOコールなしで上がることを許可
	// "custom1"                       : ユーザーが好きにカスタマイズできるルール。
	// "custom2"                       : ユーザーが好きにカスタマイズできるルール。追加してもよい。
	RuleToApply string  `json:"rule"         comment:"適用するルール。"default:"official_new"`
	CustomRules UnoRule `json:"custom-rules" comment:"ユーザーがルールをカスタマイズできる。"`
	Port        int     `json:"port"         comment:"サーバーが使用するポート番号" default:"5000"`
	LogPath     string  `json:"log"          comment:"ログファイル出力先を指定する。この指定がないときログ出力しない。" default:""`
	Silent      bool    `json:"quiet"        comment:"標準出力と標準エラー出力に出力しない。" default:"false"`
}

// 空の設定ファイルを生成する。
func CreateEmptyHjson() string {
	p := "./uno_setting_empty.hjson"
	var sb string
	// 空の設定ファイルを生成する。
	// 空の構造体から、jsonを生成し、hjsonに変換する。
	s := NewSetting()
	s.CustomRules.SetRule("official_new")
	if !true {
		// hjsonはOrderedMapSSをMarshalできないのでjsonとしてMarshalしてHjsonに変換する。
		jsonData, err := json.Marshal(*s)
		dst := &bytes.Buffer{}
		if err := json.Compact(dst, jsonData); err != nil {
			panic(err)
		}
		jsonData = dst.Bytes()
		b, err := ConvertJSONdataToHJSON(jsonData)
		if err != nil {
			panic(errors.Errorf("%v", err))
		}
		sb = string(b)
	} else {
		b, err := hjson.Marshal(*s)
		if err != nil {
			panic(errors.Errorf("%v", err))
		}
		sb = string(b)
	}

	sbs := strings.Split(sb, "\n")

	for i, l := range sbs {
		sbs[i] = strings.TrimLeftFunc(l, unicode.IsSpace) // インデントを削除する。
	}

	sb = strings.Join(sbs, "\n")
	sb = sb[2:len(sb)-2] + "\n" // 最初と最後の{}を削除する。
	sb = RemoveEmptyLines(sb)

	sb = fmt.Sprintf(`##################################################################################################################################
#
# 設定ファイルの雛形。
#
# - このファイルの名前は%vに変更すること。他のパスや名前の場合、オプション'-s'で指定すること。
#
# - hjson形式。hjson は簡単に言えばコメントが書ける書きやすいjson。hjsonの書式の詳細は公式サイトを参照。 https://hjson.github.io/try.html
#     - ファイルパス文字列は\は/で記述すること。
#     - コメントには #と//と/**/の形式が使える。
#         - ファイルサーバー上のパスなどで、//で始まる場合、コメント扱いされるのを避けるために "" で括る必要がある。
#     - 文字列は上記の//を除いて引用符で括る必要はなく、そのまま書いてよい。
#         - ヒアドキュメントのように、文字列を複数行書きたいときは以下のように'''で括ることで記述できる。下記は例。
#           '''
#           ここに複数行記述する。
#           ここに複数行記述する。
#           '''
#
# - 不要な設定値は削除してよい。
# - 設定値がファイルパスの時空白があると誤動作する可能性がある。
# - 実際にはjsonに変換されるものと思ってよい。 サブコマンド'convertwを使用して、setting.hjsonをjsonに変換できる。
# - より具体的な例が見たいときは、 ./create_installer_iso.exe template のようにして追加の雛形を生成できる。
#   また、テスト用コマンド ./create_installer_iso.exe test を実行すると更に具体的な実行例が生成できるので参考に。
#
##################################################################################################################################

%v
`, SettingFileName, sb)
	sb = IndentOutermostBraces(sb)

	sb += `
##################################################################################################################################
## "rule" で指定できるルールと説明。
#`
	// ルール一覧を出力する。
	rules := strings.Split(UnoRuleList.Print(), "\n")
	for _, r := range rules {
		sb += fmt.Sprintf("\n# %v", r)
	}
	/*
			// カスタマイズできるルール一覧。
			sb += `
		##################################################################################################################################
		# "custom-rules" で指定できるルールと説明。
		#`
			rules = strings.Split(UnoRule{}.PrintForComment(), "\n")
			for _, r := range rules {
				sb += fmt.Sprintf("\n# %v", r)
			}
	*/
	sb += `
##################################################################################################################################
`

	//log.Printf(sb)
	WriteText(&p, &sb)
	return sb
}

func NewSetting() *Setting {
	return &Setting{
		LogPath:     "",
		Silent:      false,
		RuleToApply: "official_new",
		Port:        5000,
	}
}

func (s *Setting) Print(setting_path string) string {
	if s == nil {
		return "s is nil"
	}
	return fmt.Sprintf(`
	setting hjson path       : %q
	log                      : %q
	`,
		AbsPath(setting_path),
		s.LogPath,
	)
}
