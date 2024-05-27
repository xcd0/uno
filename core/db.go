package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS PlayerInfo (
		Name TEXT,
		GUID TEXT PRIMARY KEY,
		Registered DATE,
		Email TEXT,
		IsNPC BOOLEAN
	);`
	_, err := db.Exec(query)
	return err
}

func InsertPlayer(db *sql.DB, player PlayerInfo) error {
	query := `
	INSERT INTO PlayerInfo (Name, GUID, Registered, Email, IsNPC) 
	VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(query, player.Name, player.PlayerID, player.Registered, player.Email, player.IsNPC)
	return err
}

func GetPlayerByName(db *sql.DB, name string) (*PlayerInfo, error) {
	query := `SELECT Name, GUID, Registered, Email, IsNPC FROM PlayerInfo WHERE Name = ?;`
	row := db.QueryRow(query, name)

	var player PlayerInfo
	if err := row.Scan(&player.Name, &player.PlayerID, &player.Registered, &player.Email, &player.IsNPC); err != nil {
		return nil, err
	}
	return &player, nil
}

func GetPlayerByGUID(db *sql.DB, guid string) (*PlayerInfo, error) {
	query := `SELECT Name, GUID, Registered, Email, IsNPC FROM PlayerInfo WHERE GUID = ?;`
	row := db.QueryRow(query, guid)

	var player PlayerInfo
	if err := row.Scan(&player.Name, &player.PlayerID, &player.Registered, &player.Email, &player.IsNPC); err != nil {
		return nil, err
	}
	return &player, nil
}

func GetPlayerByEmail(db *sql.DB, email string) (*PlayerInfo, error) {
	query := `SELECT Name, GUID, Registered, Email, IsNPC FROM PlayerInfo WHERE Email = ?;`
	row := db.QueryRow(query, email)

	var player PlayerInfo
	if err := row.Scan(&player.Name, &player.PlayerID, &player.Registered, &player.Email, &player.IsNPC); err != nil {
		return nil, err
	}
	return &player, nil
}

func OpenOrCreateDatabase(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Database does not exist, creating a new one.")
	} else {
		fmt.Println("Database exists, opening.")
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := InitializeDatabase(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func db_example() {

	db, err := OpenOrCreateDatabase(UnoDBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the database
	if err := InitializeDatabase(db); err != nil {
		log.Fatal(err)
	}

	// Insert a sample player
	player := PlayerInfo{
		Name:       "John Doe",
		PlayerID:   "1234-5678-9012",
		Registered: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		Email:      "john.doe@example.com",
		IsNPC:      false,
	}

	if err := InsertPlayer(db, player); err != nil {
		log.Fatal(err)
	}

	fmt.Println("PlayerInfo inserted successfully!")

	retrievedPlayer, err := GetPlayerByName(db, "John Doe")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved PlayerInfo by Name: %+v\n", retrievedPlayer)

	retrievedPlayer, err = GetPlayerByGUID(db, "1234-5678-9012")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved PlayerInfo by GUID: %+v\n", retrievedPlayer)

	retrievedPlayer, err = GetPlayerByEmail(db, "john.doe@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved PlayerInfo by Email: %+v\n", retrievedPlayer)
}
