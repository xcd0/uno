MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILDDIR     := $(MAKEFILE_DIR)
BIN          := uno
VERSION      := 0.0.1
REVISION     := `git rev-parse --short HEAD`
FLAG         := -ldflags='-X main.version=$(VERSION) -X main.revision='$(REVISION)' -s -w -extldflags="-static" -buildid=' -a -tags netgo -installsuffix -trimpath

# Detect OS
ifeq ($(OS),Windows_NT)
	BIN := $(BIN).exe
endif

all: build
build:
	@echo "Building..."
	go build -C cmd/uno -o $(BUILDDIR)$(BIN)
	@echo "Built successfully."
release:
	@echo "Building..."
	go build -C cmd/uno -o $(BUILDDIR)$(BIN) $(FLAG)
	@echo "Built successfully."
	upx --lzma $(BUILDDIR)$(BIN)

# PHONY targets
.PHONY: all server client clean release

test-register:
	make build
	`$(BUILDDIR)/uno$(EXE) server > test.log 2>&1 & pid=$$! && sleep 1 && \
		curl -s -X POST http://localhost:5000/api/players/register \
		-H "Content-Type: application/json" -d '{"name": "John"}' 2>&1 > resp.log; \
		kill $$pid`
	cat resp.log | jq "." >> test.log
	rm resp.log
	cat test.log

test-run-server:
	make build
	$(BUILDDIR)/uno$(EXE) server > test_server.log 2>&1 & echo "PID: $$!" > test_pid.log;
	ps aux | grep `cat test_pid.log | grep "PID" | awk '{print $$2}'` | grep -v "grep"
test-stop-server:
	-ps aux | grep `cat test_pid.log | grep "PID" | awk '{print $$2}'` | grep -v "grep"
	-kill `cat test_pid.log | grep "PID" | awk '{print $$2}'`
	-ps aux | grep `cat test_pid.log | grep "PID" | awk '{print $$2}'` | grep -v "grep"
	-kill `ps aux | grep uno | awk '{print $$2}'` >/dev/null 2>&1
	-ps aux | grep uno | grep -v "grep"

test-newgame:
	-time curl -s -X POST http://localhost:5000/api/game/new \
		--max-time 10 \
		-H "Content-Type: application/json" -d '{"name":"John","player_id":"518a8a2f-e690-45a8-914c-6c0ccb43d02a"}' 2>&1 > resp.log; \
	-cat resp.log | jq "." > test.log
	-rm resp.log
	-cat test.log

test-gamestate:
	make build
	$(BUILDDIR)/uno$(EXE) server &; \
		pid=$!; \
		curl -X GET http://localhost:5000/state -i; \
		kill $(pid)
test-play:
	make build
	$(BUILDDIR)/uno$(EXE) server &; \
		pid=$!; \
		curl -X GET http://localhost:5000/play -i; \
		kill $(pid)

tmp:

	#router.HandleFunc("/api/players/register", HandleRegisterPlayer).Methods("POST")
	#router.HandleFunc("/api/players/{playerId}", HandleGetPlayerInfo).Methods("GET")
	#router.HandleFunc("/api/game/new", HandleNewGame).Methods("POST")     // 新しいゲームを開始し、セッションIDを発行します。
	#router.HandleFunc("/api/game/{sessionId}/state", HandleGameState).Methods("GET")  // 指定されたセッションIDのゲーム状態を取得します。現在の状態JsonGameStateをクライアントに送信するためのGETリクエストを処理します。
	#router.HandleFunc("/api/game/{sessionId}/play", HandleClientPlay).Methods("POST") // 指定されたセッションIDのゲームにおいて、プレーヤーのアクションを処理します。クライアントがJsonClientPlayを送信するPOSTリクエストを処理します。
	#router.HandleFunc("/api/game/{sessionId}/cards", HandleCards).Methods("GET")      // ゲームで使用されるすべてのカードの詳細情報を取得します。すべてのカードのリストを送信するためのGETリクエストを処理します。
