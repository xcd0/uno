package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xcd0/uno/core"
)

// TestHandleGameState は /gamestate エンドポイントのGETリクエストをテストします。
func TestHandleGameState(t *testing.T) {
	req, err := http.NewRequest("GET", "/gamestate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGameState)

	handler.ServeHTTP(rr, req)

	// ステータスコードのチェック
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// レスポンスボディのチェック
	var gameState core.JsonGameState
	if err := json.NewDecoder(rr.Body).Decode(&gameState); err != nil {
		t.Fatal(err)
	}

	// ここで gameState の内容を具体的にチェックすることが可能です。
	// 例えば、CurrentCard が 1 であることを確認します。
	if gameState.CurrentCard != 1 {
		t.Errorf("expected CurrentCard 1, got %d", gameState.CurrentCard)
	}
}

// TestHandleClientPlay は /play エンドポイントのPOSTリクエストをテストします。
func TestHandleClientPlay(t *testing.T) {
	// テスト用の JSON データ
	playData := core.JsonClientPlay{
		PlayerId: 1,
		CardId:   101,
		Action:   "draw",
	}
	jsonData, err := json.Marshal(playData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/play", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleClientPlay)

	handler.ServeHTTP(rr, req)

	// ステータスコードのチェック
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// ここでPOSTリクエストの結果をさらに詳細にチェックすることができます。
}
