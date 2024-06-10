package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var (
	wrapperStdout io.Writer
	wrapperStderr io.Writer
)

// ヘルパー関数で関数名と行番号を取得し、フォーマットして返す
func LogCallerInfo() {
	pc, file, line, ok := runtime.Caller(1) // Caller(1) は呼び出し元の情報を取得
	if !ok {
		//return "情報を取得できませんでした"
		return
	}
	funcName := runtime.FuncForPC(pc).Name() // PC値から関数名を取得
	funcName = path.Base(funcName)           // フルパスからベース名のみ取得
	fileName := path.Base(file)              // フルパスからファイル名のみ取得
	//str := fmt.Sprintf("%s:%d %s", fileName, line, funcName)
	//fmt.Printf("%s %s:%d %s\n", time.Now().Format("15:04:05"), fileName, line, funcName)
	fmt.Printf("%s %s:%d \n", time.Now().Format("15:04:05"), fileName, line)

}

func LoggingSettings(logpath string, s *Setting) {
	logpath = func() string {
		if len(logpath) != 0 {
			return logpath
		} else if s != nil && len(s.LogPath) != 0 {
			return s.LogPath
		}
		return ""
		//return filepath.ToSlash(filepath.Join(GetCurrentDir(), fmt.Sprintf("%v.log", getFileNameWithoutExt(filepath.Base(os.Args[0])))))
	}()
	if len(logpath) != 0 {
		//log.Printf("log filepath : %v", logpath)
		logfile, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(errors.Errorf("%v", err))
		}
		wrapperStdout, wrapperStderr = io.MultiWriter(os.Stdout, logfile), io.MultiWriter(os.Stderr, logfile) // 出力をログファイル、標準出力、標準エラー出力に出力する。
	} else {
		wrapperStdout, wrapperStderr = os.Stdout, os.Stderr // 出力を標準出力、標準エラー出力に出力する。
	}
	log.SetOutput(wrapperStdout)             // logの出力先
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する
}
