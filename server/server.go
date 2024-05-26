package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/xcd0/uno/core"
)

/*
# 未実装項目

- サーバー
	- クライアントの認証。
		- 簡単な成り済まし防止の実装。
		- 別にチート対策をしたい訳ではない。httpsまではやりたくない。
	- 通信間のエラーハンドリング。
		- これは追々。
	- 複数のゲームセッション同時管理機能。
		- とりあえず無くてもいいかな。
		- 作る場合、URL/セッションID/エンドポイント/でいいかな。
	- リアルタイム通信。
		- UNO!のコールはリアルタイム性が必要。
		- 追々実装してもよい。
- クライアント
	- クライアントはコマンドライン版、webブラウザ(PWAアプリ)版の2種類作りたい。
	- UI
		- ミニマルUI→まあまあUI→見やすいUIの3種類作りたい。
		- とりあえずは最小限UI。
	- エラーハンドリング
		- ユーザーの無効入力対応。
			- とりあえず無視でいい。せいぜい無効な入力であることを通知する程度でよい。
		- 通信間のエラーハンドリング。
			- これは追々。
	- 多言語対応。
		- 追々。
*/

func UnoServer(port int, rule core.UnoRule) {

	// handleCardsはすべてのカードのリストを送信するためのGETリクエストを処理します。
	// ruleを引数渡しできないため、クロージャとしてハンドラーを定義する。
	handleCards := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		cards := (&core.JsonCardInfo{}).Init(rule)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}

	http.HandleFunc("/gamestate", HandleGameState) // HandleGameStateは現在の状態JsonGameStateをクライアントに送信するためのGETリクエストを処理します。
	http.HandleFunc("/play", HandleClientPlay)     // HandleClientPlayは、クライアントがJsonClientPlayを送信するPOSTリクエストを処理します。
	http.HandleFunc("/cards", handleCards)         // HandleCardsはすべてのカードのリストを送信するためのGETリクエストを処理します。
	log.Printf("Server starting on port %v...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func GenerateJsonGameStateForNewTurn() core.JsonGameState {
	// 未実装
	return core.JsonGameState{
		CurrentCard:      1,
		DrawnCardID:      []int{102, 205},
		Hand:             []int{101, 202, 303},
		TeammateHand:     []int{104, 208},
		YourName:         "Player1",
		OtherPlayerNames: []string{"Player2", "Player3", "Player4"},
		OtherHandCounts:  []int{3, 4, 5},
		UnoCalled:        []bool{false, true, false},
		Direction:        1,
	}
}

// HandleGameStateは現在の状態をクライアントに送信するためのGETリクエストを処理します。
func HandleGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	gameState := GenerateJsonGameStateForNewTurn()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

// HandleClientPlayは、クライアントがプレイを送信するPOSTリクエストを処理します。
func HandleClientPlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var clientPlay core.JsonClientPlay
	err := json.NewDecoder(r.Body).Decode(&clientPlay)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	log.Printf("Received play: %+v\n", clientPlay)

	// ここでプレイを処理し、新しいゲームの状態や結果を送り返すことができます。

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Send back that everything was fine
}
