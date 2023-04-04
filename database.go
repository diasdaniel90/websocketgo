package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func saveToDatabase(db *sql.DB, p *Payload) error {
	// Consultar se já existe registro com o id_bet informado
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT id_bet FROM api_serverresult WHERE id_bet = ?)", p.IdBet).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	if exists {
		return nil
	}

	// Preparar o command SQL para inserir dados na tabela
	stmt, err := db.Prepare("INSERT INTO  api_serverresult (ID_bet, bet_color, bet_roll, `timestamp`, bet_status, total_red_eur_bet, total_red_bets_placed, total_white_eur_bet, total_white_bets_placed, total_black_eur_bet, total_black_bets_placed, total_bets_placed, total_eur_bet, total_retention_eur) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Executar o command SQL com valores
	_, err = stmt.Exec(p.IdBet, p.Color, p.Roll, p.Timestamp, p.Status, p.TotalRedEurBet, p.TotalRedBetsPlaced, p.TotalWhiteEurBet, p.TotalWhiteBetsPlaced, p.TotalBlackEurBet, p.TotalBlackBetsPlaced, p.TotalBetsPlaced, p.TotalEurBet, p.TotalRetentionEur)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Dados inseridos com sucesso!")

	return nil
}

func saveToDatabaseUsers(db *sql.DB, p *Payload) error {
	for _, bet := range (*p).Bets {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT ID_bet_uniqa FROM api_userresult WHERE ID_bet_uniqa = ?)", bet.IDBetUser).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error connecting to websocket: %w", err)
		}
		if exists {
			//log.Println("registro já existe")
			return nil
		}
		tBetUser, _ := time.Parse(layout, p.CreatedAt)
		p.Timestamp = tBetUser.Unix()
		stmt, err := db.Prepare("INSERT INTO api_userresult (ID_bet, ID_bet_uniqa, `timestamp`, color, amount, currency_type,user) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(p.IdBet, bet.IDBetUser, p.Timestamp, bet.Color, bet.Amount, bet.CurrencyType, bet.User.IDStr)
		if err != nil {
			panic(err.Error())
		}
	}

	log.Println("Dados inseridos com sucesso!")

	return nil
}
