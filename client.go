package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// FetchGameState はサーバーから現在のゲーム状態を取得します。
func FetchGameState(url string) (*JsonGameState, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gameState JsonGameState
	err = json.NewDecoder(resp.Body).Decode(&gameState)
	if err != nil {
		return nil, err
	}
	return &gameState, nil
}

// SendClientAction はクライアントの行動をサーバーに送信します。
func SendClientAction(url string, action *JsonClientPlay) error {
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Server Response:", string(body))
	return nil
}

func GenerateJsonClientPlay() JsonClientPlay {
	// 未実装
	return JsonClientPlay{
		Discard:           []int{101}, // 例として101番のカードを捨てる
		CallUno:           true,
		CallUnoOnPlayerID: -1,
	}
}

func UnoClient(port int) {
	// サーバーからゲーム状態を取得
	gameStateUrl := fmt.Sprintf("http://localhost:%v/gamestate", port)
	gameState, err := FetchGameState(gameStateUrl)
	if err != nil {
		log.Fatalf("Failed to fetch game state: %v", err)
	}
	fmt.Printf("Current Game State: %+v\n", gameState)

	// サーバーにクライアントの行動を送信
	playUrl := fmt.Sprintf("http://localhost:%v/play", port)

	action := GenerateJsonClientPlay()
	err = SendClientAction(playUrl, &action)
	if err != nil {
		log.Fatalf("Failed to send client action: %v", err)
	}
}
