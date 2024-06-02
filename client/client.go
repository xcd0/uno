package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xcd0/uno/core"
)

type ClientState int

const (
	Start ClientState = iota
	Select
	Single
	Multi
	End
	Exit
)

func UnoClientStart() (next ClientState) {
	log.Printf("UnoClientStart: start")
	next = Select
	// なんかこういい感じのスタート画面を出す。

	box := tview.NewBox().
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("A [red]c[yellow]o[green]l[darkcyan]o[blue]r[darkmagenta]f[red]u[yellow]l[white] [black:red]c[:yellow]o[:green]l[:darkcyan]o[:blue]r[:darkmagenta]f[:red]u[:yellow]l[white:-] [::bu]title")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}

	app := tview.NewApplication()
	button := app.NewButton("Hit Enter to close").SetSelectedFunc(func() {
		app.Stop()
	})
	button.SetBorder(true).SetRect(0, 0, 22, 3)

	if err := app.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
	log.Printf("UnoClientStart: end (next-> %v)", next)
	return
}
func UnoClientSelect() (next ClientState) {
	log.Printf("%v: start", core.FuncName())
	// ユーザーと遊ぶモードを選ぶ。

	if true {
		next = Single
	} else {
		next = Multi
	}
	log.Printf("%v: end (next-> %v)", core.FuncName(), next)
	return
}
func UnoClientPlaySingle() (next ClientState) {
	log.Printf("%v: start", core.FuncName())
	next = End
	log.Printf("%v: end (next-> %v)", core.FuncName(), next)
	return
}
func UnoClientPlayMulti() (next ClientState) {
	log.Printf("%v: start", core.FuncName())
	next = End
	log.Printf("%v: end (next-> %v)", core.FuncName(), next)
	return
}
func UnoClientEnd() (next ClientState) {
	log.Printf("%v: start", core.FuncName())
	// 終了時の表示をする。
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
				next = Exit
			} else if buttonLabel == "Cancel" {
				app.Stop()
				next = Select
			}
		})
	if err := app.SetRoot(modal, false).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	log.Printf("%v: end (next-> %v)", core.FuncName(), next)
	return
}

func UnoClient(port int) {

	var next ClientState = Start

	for {
		switch next {
		case Start:
			next = UnoClientStart()
		case Select:
			next = UnoClientSelect()
		case Single:
			next = UnoClientPlaySingle()
		case Multi:
			next = UnoClientPlayMulti()
		case End:
			next = UnoClientEnd()
		case Exit:
			return
		}
	}

	return

	// サーバーからゲーム状態を取得
	gameStateUrl := fmt.Sprintf("http://localhost:%v/gamestate", port)
	var gameState *core.JsonGameState
	for i := 10; i > 0; i-- {
		var err error
		gameState, err = FetchGameState(gameStateUrl)
		if err != nil {
			log.Printf("Failed to fetch game state: %v", err)
			time.Sleep(1 * time.Second)
		}
		if i == 1 {
			log.Fatalf("Failed to fetch game state: %v", err)
		}
	}
	fmt.Printf("Current Game State: %+v\n", gameState)

	// サーバーにクライアントの行動を送信
	playUrl := fmt.Sprintf("http://localhost:%v/play", port)

	action := GenerateJsonClientPlay()
	if err := SendClientAction(playUrl, &action); err != nil {
		log.Fatalf("Failed to send client action: %v", err)
	}
}

// FetchGameState はサーバーから現在のゲーム状態を取得します。
func FetchGameState(url string) (*core.JsonGameState, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gameState core.JsonGameState
	err = json.NewDecoder(resp.Body).Decode(&gameState)
	if err != nil {
		return nil, err
	}
	return &gameState, nil
}

// SendClientAction はクライアントの行動をサーバーに送信します。
func SendClientAction(url string, action *core.JsonClientPlay) error {
	jsonData, err := json.Marshal(action)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// レスポンスの確認
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Server Response:", string(body))
	return nil
}

func GenerateJsonClientPlay() core.JsonClientPlay {
	// 未実装
	return core.JsonClientPlay{
		Discard:           []int{101}, // 例として101番のカードを捨てる
		CallUno:           true,
		CallUnoOnPlayerID: -1,
	}
}
