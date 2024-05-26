package core

import (
	"fmt"
	"path"
	"runtime"
	"time"
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
