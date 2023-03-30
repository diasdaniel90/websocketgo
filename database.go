package main

import (
	"database/sql"
	"log"
)

func saveToDatabase(db *sql.DB, p *Payload) error {
	log.Print("se vira para salvar essa parada no banco de dados", p)

	// Consultar se já existe registro com o id_bet informado
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT id_bet FROM source_double.api_serverresult WHERE id_bet = ?)", p.IdBet).Scan(&exists)
	if err != nil {
		log.Println(err)
		return err
	}
	if exists {
		log.Println("registro já existe")
		return nil
	}

	// Preparar o comando SQL para inserir dados na tabela "usuarios"
	stmt, err := db.Prepare("INSERT INTO  source_double.api_serverresult (ID_bet, bet_color, bet_roll, `timestamp`, bet_status, total_red_eur_bet, total_red_bets_placed, total_white_eur_bet, total_white_bets_placed, total_black_eur_bet, total_black_bets_placed, total_bets_placed, total_eur_bet, total_retention_eur) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Executar o comando SQL com valores para nome e idade
	_, err = stmt.Exec(p.IdBet, p.Color, p.Roll, p.Timestamp, p.Status, p.TotalRedEurBet, p.TotalRedBetsPlaced, p.TotalWhiteEurBet, p.TotalWhiteBetsPlaced, p.TotalBlackEurBet, p.TotalBlackBetsPlaced, p.TotalBetsPlaced, p.TotalEurBet, p.TotalRetentionEur)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Dados inseridos com sucesso!")

	return nil
}

func saveToDatabaseUsers(db *sql.DB, p *Payload) error {
	log.Print("se vira para salvar essa parada no banco de dados", p)

	// Consultar se já existe registro com o id_bet informado
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT ID_bet_uniqa FROM source_double.api_userresult WHERE ID_bet_uniqa = ?)", p.IdBet).Scan(&exists)
	if err != nil {
		log.Println(err)
		return err
	}
	if exists {
		log.Println("registro já existe")
		return nil
	}

	// Preparar o comando SQL para inserir dados na tabela "usuarios"
	stmt, err := db.Prepare("INSERT INTO source_double.api_userresult (ID_bet, bet_color, bet_roll, `timestamp`, bet_status, total_red_eur_bet, total_red_bets_placed, total_white_eur_bet, total_white_bets_placed, total_black_eur_bet, total_black_bets_placed, total_bets_placed, total_eur_bet, total_retention_eur) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Executar o comando SQL com valores para nome e idade
	_, err = stmt.Exec(p.IdBet, p.Color, p.Roll, p.Timestamp, p.Status, p.TotalRedEurBet, p.TotalRedBetsPlaced, p.TotalWhiteEurBet, p.TotalWhiteBetsPlaced, p.TotalBlackEurBet, p.TotalBlackBetsPlaced, p.TotalBetsPlaced, p.TotalEurBet, p.TotalRetentionEur)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Dados inseridos com sucesso!")

	return nil
}
