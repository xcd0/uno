package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/xcd0/uno/core"

	_ "github.com/mattn/go-sqlite3"
)

type SessionInfo struct {
	State      *core.State
	PlayerUID  []string
	PlayerName map[string]string // UID -> Name
}

// 仮のゲームセッションとプレイヤー管理
var gameSessions = struct {
	sync.Mutex
	sessions map[string]*SessionInfo
}{sessions: make(map[string]*SessionInfo)}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketコネクション管理
var wsClients = struct {
	sync.Mutex
	connections map[string]*websocket.Conn
}{connections: make(map[string]*websocket.Conn)}

// プレイヤー情報を管理するためのマップ
var playerStore = struct {
	sync.Mutex
	players map[string]core.PlayerInfo
	names   map[string]string // 名前のユニーク性を保持するためのマップ
}{players: make(map[string]core.PlayerInfo), names: make(map[string]string)}

var (
	HandleNewGame    func(w http.ResponseWriter, r *http.Request)
	HandleGameState  func(w http.ResponseWriter, r *http.Request)
	HandleClientPlay func(w http.ResponseWriter, r *http.Request)
	HandleCards      func(w http.ResponseWriter, r *http.Request)
	HandleWebSocket  func(w http.ResponseWriter, r *http.Request)
)

func UnoServer(state *core.State, wg *sync.WaitGroup) {
	log.Printf("UnoServer: wake up")
	wg.Add(1)
	defer wg.Done()
	rule := state.Rule
	port := state.Setting.Port
	//yourname := "You"
	cards := (&core.JsonCardInfo{}).Init(&rule)

	state.Players = []string{}

	// DBを開く。
	db, err := core.OpenOrCreateDatabase(core.UnoServerDBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := core.InitializeDatabase(db); err != nil {
		log.Fatal(err)
	}

	//log.Printf("%v", state.Setting.Print())

	HandleNewGame = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var newPlayer struct {
			Name     string `json:"name"`
			PlayerID string `json:"player_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&newPlayer); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		sessionID := uuid.New().String()
		gameSessions.Lock()
		gameSessions.sessions[sessionID] = &SessionInfo{
			State:      state,
			PlayerUID:  []string{newPlayer.PlayerID},
			PlayerName: map[string]string{},
		}

		var si *SessionInfo = gameSessions.sessions[sessionID]

		si.PlayerName[newPlayer.PlayerID] = newPlayer.Name
		//si.PlayerUID = append(si.PlayerUID, newPlayer.PlayerID)
		gameSessions.Unlock()

		log.Printf("新しいセッションを開始しました: %s", sessionID)

		playerAdded := make(chan struct{})
		allPlayersAdded := make(chan struct{})

		// プレイヤーが集まるのを待つ
		go waitForPlayers(sessionID, 4, playerAdded, allPlayersAdded)

		// セッションが開始されてから5秒間待機して参加者がいない場合はNPCを追加する
		select {
		case <-time.After(5 * time.Second):
			gameSessions.Lock()
			log.Printf("5秒経過したため、NPCを追加します。")

			for i := len(si.State.Players); i <= 3; i++ {
				si.State.Players = append(si.State.Players, fmt.Sprintf("NPC%d", i))
				puid := uuid.New().String()
				si.PlayerUID = append(si.PlayerUID, puid)
				si.PlayerName[puid] = si.State.Players[i]
			}
			gameSessions.Unlock()
			close(playerAdded)
		case <-allPlayersAdded:
			log.Printf("プレイヤーが集まりました: %s", sessionID)
		}

		log.Printf("session id: %v, players: %v", sessionID, si.State.Players)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"sessionId": sessionID, "message": "New game session started successfully."})
	}

	// HandleGameStateは現在の状態をクライアントに送信するためのGETリクエストを処理します。
	HandleGameState = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		sessionID := r.URL.Query().Get("sessionId")
		gameSessions.Lock()
		var si *SessionInfo = gameSessions.sessions[sessionID]
		//session, exists := gameSessions.sessions[sessionID]
		_, exists := gameSessions.sessions[sessionID]
		gameSessions.Unlock()
		if !exists {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		// ここで適切なゲーム状態を設定する必要があります
		gameState := core.JsonGameState{
			CurrentCard:      si.State.LastCard.Type,
			DrawnCardID:      si.State.DrawnCardID,
			Hand:             si.State.HandID(),
			TeammateHand:     []int{}, //session.TeammateHand,
			YourName:         si.State.Name(),
			OtherPlayerNames: si.State.OtherPlayerNames(),
			OtherHandCounts:  si.State.OtherHandCounts(),
			UnoCalled:        si.State.UnoCalled,
			Direction:        si.State.Turn,
			RuleToApply:      si.State.Setting.RuleToApply,
			CustomRules:      si.State.Rule,
			Port:             si.State.Setting.Port,
		}
		//gameState := GenerateJsonGameStateForNewTurn(yourname, state)

		//func (jgs *JsonGameState) Init(yourname string, state *State) JsonGameState {

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gameState)
	}

	// HandleClientPlayは、クライアントがプレイを送信するPOSTリクエストを処理します。
	HandleClientPlay = func(w http.ResponseWriter, r *http.Request) {
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
		sessionID := r.URL.Query().Get("sessionId")
		gameSessions.Lock()
		var si *SessionInfo = gameSessions.sessions[sessionID]
		//session, exists := gameSessions.sessions[sessionID]
		_, exists := gameSessions.sessions[sessionID]
		gameSessions.Unlock()
		if !exists {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}

		log.Printf("Received play: %+v\n", clientPlay)
		// プレイされたカードの更新

		si.State.LastCard = &cards.CardAll[clientPlay.Discard[len(clientPlay.Discard)-1]]

		// 次のプレイヤーのターンを通知する処理
		nextPlayerID := si.State.NextPlayerID()
		notifyPlayerTurn(sessionID, nextPlayerID)

		// ここでプレイを処理し、新しいゲームの状態や結果を送り返すことができます。
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Send back that everything was fine
	}

	// HandleCardsはすべてのカードのリストを送信するためのGETリクエストを処理します。
	// ruleを引数渡しできないため、クロージャとしてハンドラーを定義する。
	HandleCards = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		//cards := (&core.JsonCardInfo{}).Init(&rule)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}

	HandleWebSocket = func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		playerID := vars["playerId"]
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not upgrade to websocket", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		// プレイヤーIDを取得
		//playerID := r.URL.Query().Get("playerId")
		wsClients.Lock()
		wsClients.connections[playerID] = conn
		wsClients.Unlock()
		for {
			// メッセージの読み取り（必要に応じて実装）
			_, _, err := conn.ReadMessage()
			if err != nil {
				wsClients.Lock()
				delete(wsClients.connections, playerID)
				wsClients.Unlock()
				break
			}
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/players/register", HandleRegisterPlayer).Methods("POST")
	router.HandleFunc("/api/players/{playerId}", HandleGetPlayerInfo).Methods("GET")
	router.HandleFunc("/api/game/new", HandleNewGame).Methods("POST")                 // 新しいゲームを開始し、セッションIDを発行します。
	router.HandleFunc("/api/game/{sessionId}/state", HandleGameState).Methods("GET")  // 指定されたセッションIDのゲーム状態を取得します。現在の状態JsonGameStateをクライアントに送信するためのGETリクエストを処理します。
	router.HandleFunc("/api/game/{sessionId}/play", HandleClientPlay).Methods("POST") // 指定されたセッションIDのゲームにおいて、プレーヤーのアクションを処理します。クライアントがJsonClientPlayを送信するPOSTリクエストを処理します。
	router.HandleFunc("/api/game/{sessionId}/cards", HandleCards).Methods("GET")      // ゲームで使用されるすべてのカードの詳細情報を取得します。すべてのカードのリストを送信するためのGETリクエストを処理します。
	router.HandleFunc("/api/ws/{playerId}", HandleWebSocket)                          // 参加中のゲームのリアルタイムイベント用webソケット。
	router.HandleFunc("/api/ws/subscribe/{playerId}", HandleWebSocket)                // プレイヤーがリアルタイムのゲームイベントを購読するためのWebSocket接続を提供します。

	// http.HandleFunc("/api/{sessionId}/game/new", HandleNewGame)     // 新しいゲームを開始し、セッションIDを発行します。
	// http.HandleFunc("/api/{sessionId}/game/state", HandleGameState) // 指定されたセッションIDのゲーム状態を取得します。現在の状態JsonGameStateをクライアントに送信するためのGETリクエストを処理します。
	// http.HandleFunc("/api/{sessionId}/game/play", HandleClientPlay) // 指定されたセッションIDのゲームにおいて、プレーヤーのアクションを処理します。クライアントがJsonClientPlayを送信するPOSTリクエストを処理します。
	// http.HandleFunc("/api/{sessionId}/game/cards", HandleCards)     // ゲームで使用されるすべてのカードの詳細情報を取得します。すべてのカードのリストを送信するためのGETリクエストを処理します。
	// http.HandleFunc("/ws", HandleWebSocket)                         // websoket用。

	//log.Printf("Server starting on port %v...", port)
	// if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
	log.Printf("UnoServer: listen localhost:%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

func notifyPlayerTurn(sessionID string, playerID int) {
	wsClients.Lock()
	defer wsClients.Unlock()

	playerUID := gameSessions.sessions[sessionID].PlayerUID[playerID]
	if conn, ok := wsClients.connections[playerUID]; ok {
		message := fmt.Sprintf("Your turn in session %s", sessionID)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Error notifying player %s: %v", playerUID, err)
			conn.Close()
			delete(wsClients.connections, playerUID)
		}
	}
}

// プレイヤーが集まるのを待つための関数
func waitForPlayers(sessionID string, requiredPlayers int, playerAdded chan struct{}, allPlayersAdded chan struct{}) {
	log.Printf("waitForPlayers: start")
	for {
		gameSessions.Lock()
		if len(gameSessions.sessions[sessionID].PlayerUID) >= requiredPlayers {
			gameSessions.Unlock()
			close(allPlayersAdded)
			return
		}
		gameSessions.Unlock()
		select {
		case <-playerAdded:
			return
		case <-time.After(100 * time.Millisecond): // ポーリング間隔
		}
	}
	log.Printf("waitForPlayers: done")
}

func GenerateJsonGameStateForNewTurn(yourname string, state *core.State) core.JsonGameState {
	gs := core.JsonGameState{}
	return gs.Init(yourname, state)
}

// プレイヤーを登録するエンドポイント
func HandleRegisterPlayer(w http.ResponseWriter, r *http.Request) {
	var newPlayer struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPlayer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 名前の重複チェック
	playerStore.Lock()
	defer playerStore.Unlock()
	if _, exists := playerStore.names[newPlayer.Name]; exists {
		http.Error(w, "Name already taken", http.StatusConflict)
		return
	}

	playerId := uuid.New().String()
	playerInfo := core.PlayerInfo{
		PlayerID:   playerId,
		Name:       newPlayer.Name,
		Registered: time.Now(),
		IsNPC:      false,
		Email:      "",
	}
	playerStore.players[playerId] = playerInfo
	playerStore.names[newPlayer.Name] = playerId

	// Open a connection to the SQLite3 database
	db, err := sql.Open("sqlite3", "./server.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the database
	if err := core.InitializeDatabase(db); err != nil {
		log.Fatal(err)
	}

	// Insert a sample player
	player := core.PlayerInfo{
		Name:       "John Doe",
		PlayerID:   "1234-5678-9012",
		Registered: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Email:      "john.doe@example.com",
		IsNPC:      false,
	}

	if err := core.InsertPlayer(db, player); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerInfo)
}

// プレイヤー情報を取得するエンドポイント
func HandleGetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId := vars["playerId"]

	playerStore.Lock()
	playerInfo, exists := playerStore.players[playerId]
	playerStore.Unlock()
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerInfo)
}
