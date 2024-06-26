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

	var dev = "EA games"
	developer_games, err := gamesByDeveloper(dev)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("games developed by %v: %v\n", dev, developer_games)

	var ID int64 = 4
	game, err := gamesByID(ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("games found with ID number %v: %v\n", ID, game)

	var Price float32 = 20
	games_in_budget, err := sortByMaxPrice(Price)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("games found under $%v: %v\n", Price, games_in_budget)

}

func gamesByDeveloper(developer string) ([]Game, error) {
	var games_slice []Game

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
		games_slice = append(games_slice, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gameByDeveloper %q: %v", developer, err)
	}
	return games_slice, nil
}

func gamesByID(id int64) (Game, error) {
	var game Game

	row := db.QueryRow("SELECT * FROM games WHERE ID = ?", id)
	if err := row.Scan(&game.ID, &game.Title, &game.Developer, &game.Price); err != nil {

		if err == sql.ErrNoRows {
			return game, fmt.Errorf("gamesById %d: no such game", id)
		}
		return game, fmt.Errorf("gamesById %d: %v", id, err)
	}
	return game, nil

}

func sortByMaxPrice(price float32) ([]Game, error) {
	var within_budget []Game

	rows, err := db.Query("SELECT * FROM games WHERE PRICE <= ?", price)
	if err != nil {
		return nil, fmt.Errorf("sortByMaxPrice %v: %v", price, err)
	}
	defer rows.Close()

	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Developer, &game.Price); err != nil {
			return nil, fmt.Errorf("sortByMaxPrice %v: %v", price, err)
		}
		within_budget = append(within_budget, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sortByMaxPrice %v: %v", price, err)
	}
	return within_budget, nil
}
