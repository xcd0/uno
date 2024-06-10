package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/pkg/errors"
	"github.com/xcd0/uno/core"
)

type Args struct {
	ArgsProto        *ArgsServer           `arg:"subcommand:proto"      help:"簡易実装。"`
	ArgsServer       *ArgsServer           `arg:"subcommand:server"     help:"サーバーモードとして起動する。"`
	ArgsClient       *ArgsClient           `arg:"subcommand:client"     help:"クライアントモードとして起動する。"`
	ArgsServerClient *ArgsServerClient     `arg:"subcommand:both"       help:"サーバーとクライアント同時に起動する。"`
	ArgsSolo         *ArgsSolo             `arg:"subcommand:solo"       help:"一人プレイモードで起動する。"`
	CreateEmptyHjson *ArgsCreateEmptyHjson `arg:"subcommand:init"       help:"空の設定ファイルを生成する。"`
	ConvertToJson    *ArgsConvert          `arg:"subcommand:convert"    help:"設定ファイルをjsonに変換する。"`
	VersionSub       *ArgsVersion          `arg:"subcommand:version"    help:"バージョン番号を出力する。-vと同じ。" `

	PlayerNames     []string `arg:"-p,--player-names"     help:"参加者の名前。"`
	NumberOfPlayers int      `arg:"-n,--number-of-player" help:"参加者全体の数。参加者の名前の数が少ない場合NPCが入る。参加者の名前の数が多い場合そちらに合わせる。" default:"4"`
	SettingPath     string   `arg:"-s,--setting"          help:"使用したい設定ファイルパス指定。指定がなければ./uno_setting.hjsonを使用する。" placeholder:"FILE"`
	LogPath         string   `arg:"-l,--logfile"          help:"ログファイル出力先を指定する。設定ファイルで指定されておらず、この指定がないときログ出力しない。" placeholder:"FILE"`
	Version         bool     `arg:"-v,--version"          help:"バージョン番号を出力する。サブコマンドversionと同じ。" `
	Help            bool
}

func (Args) Description() string {
	// TODO: 実装
	return fmt.Sprintf(`%v version %v.%v

	Description:
	なんかこういい感じの説明
	`,
		filepath.Base(os.Args[0]), version, revision)

}

// ShowHelp() で使う
var parser *arg.Parser

func ShowHelp() {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	fmt.Printf("%v\n", buf.String())
	os.Exit(1)
}

// 引数解析
func ArgParse() *Args {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する
	ThisProgramPath = core.AbsPath(os.Args[0])

	args := &Args{
		SettingPath: fmt.Sprintf("./%v", core.SettingFileName),
		PlayerNames: []string{"You", "NPC1", "NPC2", "NPC3"},
	}

	{
		var err error
		parser, err = arg.NewParser(
			arg.Config{
				Program: filepath.Base(filepath.ToSlash(os.Args[0])),
			},
			args)
		if err != nil {
			panic(errors.Errorf("%v", err))
		}
		if len(os.Args) == 1 {
			args.Help = true
			ShowHelp()
		} else if err = parser.Parse(os.Args[1:]); err != nil {
			panic(errors.Errorf("%v", err))
		}
		// --versionがなぜかtrueにならないので仕方なくチェック
		for _, arg := range os.Args[1:] {
			if arg == "--version" {
				args.Version = true
				break
			}
			if arg == "--help" || arg == "-h" {
				args.Help = true
				ShowHelp()
			}
		}
		//log.Printf("%v", args.Print())
	}

	if args.Version || args.VersionSub != nil {
		fmt.Printf("%v version %v.%v\n", filepath.Base(os.Args[0]), version, revision)
		os.Exit(1)
	}

	// 引数として必須な何れかが欠けている場合ヘルプを出力して終了する。
	if true && //
		args.ArgsProto == nil &&
		args.ArgsServer == nil &&
		args.ArgsClient == nil &&
		args.ArgsServerClient == nil &&
		args.ArgsSolo == nil &&
		args.CreateEmptyHjson == nil &&
		args.ConvertToJson == nil &&
		// args.Readme == false &&
		true {
		ShowHelp() // go-argsの生成するヘルプ文字列を取得して出力する。
	}
	return args
}

func (args *Args) Print() {
	log.Printf(`
	ArgsServer       : %5v : サーバーモードとして起動する。
	ArgsClient       : %5v : クライアントモードとして起動する。
	ArgsServerClient : %5v : サーバーとクライアント同時に起動する。
	ArgsSolo         : %5v : 一人プレイモードで起動する。
	CreateEmptyHjson : %5v : 空の設定ファイルを生成する。
	ConvertToJson    : %5v : 設定ファイルをjsonに変換する。
	VersionSub       : %5v : バージョン番号を出力する。-vと同じ。
	PlayerNames      : %5v : 参加者の名前。
	NumberOfPlayers  : %5v : 参加者全体の数。参加者の名前の数が少ない場合NPCが入る。参加者の名前の数が多い場合そちらに合わせる。
	SettingPath      : %5v : 使用したい設定ファイルパス指定。指定がなければ./uno_setting.hjsonを使用する。
	LogPath          : %5v : ログファイル出力先を指定する。設定ファイルで指定されておらず、この指定がないときログ出力しない。
	Version          : %5v : バージョン番号を出力する。サブコマンドversionと同じ。
	Help             : %5v : ヘルプを出力する。
	`,
		args.ArgsServer,
		args.ArgsClient,
		args.ArgsServerClient,
		args.ArgsSolo,
		args.CreateEmptyHjson,
		args.ConvertToJson,
		args.VersionSub,
		args.PlayerNames,
		args.NumberOfPlayers,
		args.SettingPath,
		args.LogPath,
		args.Version,
		args.Help,
	)
}

type ArgsVersion struct{}
type ArgsConvert struct {
	//Tab    bool   `arg:"--tab"   help:"jsonを見やすくformatして出力する。インデントにtabを使用する。"`
	Mini   bool   `arg:"-m,--mini" help:"minifyしたjsonを出力する。"`
	Space  int    `arg:"--space"   help:"jsonを見やすくformatして出力する。インデントに指定個数の半角空白を使用する。負の数の時無視される。" default:"-1"`
	Output string `arg:"--output"  help:"指定のパスにテキストファイルとして出力する。"`
	Color  bool   `arg:"--color" help:"色を付ける。そのあとにパイプで処理するとバグる。"`
}
type ArgsCreateEmptyHjson struct {
}
type ArgsServer struct { // サーバーモードとして起動する。
}
type ArgsClient struct { // クライアントモードとして起動する。
}
type ArgsServerClient struct { // サーバーとクライアント同時に起動する。
}
type ArgsSolo struct { // 一人プレイモードで起動する。
}
