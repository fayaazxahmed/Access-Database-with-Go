package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Game struct {
	ID        int64
	Title     string
	Developer string
	Price     float32
}

func main() {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "games",
		AllowNativePasswords: true,
	}
	fmt.Println(cfg.User, cfg.Passwd)
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	games, err := gamesByDeveloper("EA games")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("games found: %v\n", games)

}

func gamesByDeveloper(developer string) ([]Game, error) {
	var games []Game

	rows, err := db.Query("SELECT * FROM games WHERE developer = ?", developer)
	if err != nil {
		return nil, fmt.Errorf("gamesByDeveloper %q: %v", developer, err)
	}
	defer rows.Close()

	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Developer, &game.Price); err != nil {
			return nil, fmt.Errorf("gamesByDeveloper %q: %v", developer, err)
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gameByDeveloper %q: %v", developer, err)
	}
	return games, nil
}
