package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xcd0/uno/core"
)

var state *core.State = func() *core.State {
	rule := &core.UnoRule{}
	state := core.NewState(
		[]string{"You", "NPC1", "NPC2", "NPC3"},
		4,
		rule.Init().SetRule(core.UnoRuleList[0].Name),
	)
	state.Setting = core.ReadSetting("./uno_setting.hjson", core.NewSetting())
	state.Setting.Port = 8080
	core.LoggingSettings("", state.Setting)
	return state
}()

func TestHandleNewGame(t *testing.T) {
	wg := &sync.WaitGroup{}

	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	newPlayer := map[string]string{"name": "testplayer", "player_id": uuid.New().String()}
	newPlayerJSON, _ := json.Marshal(newPlayer)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%v/api/game/new", state.Setting.Port), bytes.NewBuffer(newPlayerJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleNewGame)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"sessionId"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleGameState(t *testing.T) {
	wg := &sync.WaitGroup{}
	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/api/game/{sessionId}/state", state.Setting.Port), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGameState)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"CurrentCard"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleClientPlay(t *testing.T) {
	wg := &sync.WaitGroup{}
	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	clientPlay := core.JsonClientPlay{
		Discard: []int{1, 2, 3},
	}
	clientPlayJSON, _ := json.Marshal(clientPlay)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%v/api/game/{sessionId}/play", state.Setting.Port), bytes.NewBuffer(clientPlayJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleClientPlay)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHandleCards(t *testing.T) {
	wg := &sync.WaitGroup{}
	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/api/game/{sessionId}/cards", state.Setting.Port), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCards)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"cards"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleRegisterPlayer(t *testing.T) {
	wg := &sync.WaitGroup{}
	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	newPlayer := map[string]string{"name": "testplayer"}
	newPlayerJSON, _ := json.Marshal(newPlayer)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%v/api/players/register", state.Setting.Port), bytes.NewBuffer(newPlayerJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleRegisterPlayer)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"PlayerID"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHandleGetPlayerInfo(t *testing.T) {
	wg := &sync.WaitGroup{}
	go UnoServer(state, wg)

	time.Sleep(1 * time.Second) // Server startup delay

	playerID := uuid.New().String()
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/api/players/%v", state.Setting.Port, playerID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGetPlayerInfo)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"PlayerID"`
	if !bytes.Contains(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
