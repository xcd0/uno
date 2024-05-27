package core

import "time"

// PlayerInfo はプレイヤー情報を保持します。
type PlayerInfo struct {
	Name       string    `json:"name"`
	PlayerID   string    `json:"player_id"`
	Registered time.Time `json:"registered"`
	Email      string    `json:"mail"`
	IsNPC      bool      `json:"is_npc"`
}
