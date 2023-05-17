package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func saveToDatabase(dbConexao *sql.DB, pload *payloadStruct) error {
	// Consultar se já existe registro com o id_bet informado
	var exists bool

	err := dbConexao.QueryRow(
		"SELECT EXISTS(SELECT id_bet FROM api_serverresult WHERE id_bet = ?)",
		pload.IDBet).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error SELECT EXISTS: %w", err)
	}

	if exists {
		log.Println("aviso registro já existe")

		return nil
	}

	// Preparar o command SQL para inserir dados na tabela
	stmt, err := dbConexao.Prepare(
		"INSERT INTO  api_serverresult" +
			"(ID_bet, bet_color, bet_roll, `timestamp`, bet_status, " +
			"total_red_eur_bet, total_red_bets_placed, total_white_eur_bet, total_white_bets_placed, " +
			"total_black_eur_bet, total_black_bets_placed, total_bets_placed, total_eur_bet, total_retention_eur) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// Executar o command SQL com valores
	_, err = stmt.Exec(
		pload.IDBet, pload.Color, pload.Roll, pload.Timestamp, pload.Status,
		pload.TotalRedEurBet, pload.TotalRedBetsPlaced, pload.TotalWhiteEurBet, pload.TotalWhiteBetsPlaced,
		pload.TotalBlackEurBet, pload.TotalBlackBetsPlaced, pload.TotalBetsPlaced, pload.TotalEurBet, pload.TotalRetentionEur)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Dados de Status inseridos com sucesso!")

	return nil
}

func saveToDatabaseUsers(dbConexao *sql.DB, pload payloadStruct) {
	stmt, err := dbConexao.Prepare(
		"INSERT IGNORE INTO api_gouserresults" +
			"(ID_bet, ID_bet_uniqa, `timestamp`, color, amount, currency_type,user) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	tBetUser, _ := time.Parse(layout, pload.CreatedAt)

	pload.Timestamp = tBetUser.Unix()

	for _, bet := range pload.Bets {
		_, err = stmt.Exec(
			pload.IDBet, bet.IDBetUser, pload.Timestamp, bet.Color, bet.Amount, bet.CurrencyType, bet.User.IDStr)
		if err != nil {
			panic(err.Error())
		}
	}
	// log.Println("registro do user inserido com sucesso", pload.IDBet)
}

func saveToDatabaseBets(dbConexao *sql.DB, betsLoad *[]betBotStruct) error {
	// log.Println("print", pload)
	// log.Printf("O tipo de pload é %T\n", pload)
	// log.Println("print", &pload)
	// log.Printf("O tipo de pload é %T\n", &pload)
	// defer dbConexao.Close()
	stmt, err := dbConexao.Prepare(
		"INSERT INTO api_gocontrolbetresult" +
			"(ID_bet, `timestamp`, timestamp_signal, color, source, win, " +
			"status, gale, amount, balanceWin) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error Prepare: %w", err)
	}
	defer stmt.Close()

	for _, value := range *betsLoad {
		// log.Println("vai gravar ", bet.IDBetUser)
		// timestamp_, _ := time.Parse(layout, betsLoad.CreatedAt)
		// betsLoad.Timestamp = timestamp_.Unix()
		log.Println(value)
		_, err = stmt.Exec(
			value.idBet, value.timestamp, value.timestampSinal, value.color, value.source, value.win,
			value.status, value.gale, value.amount, value.balanceWin)

		if err != nil {
			return fmt.Errorf("error Exec: %w", err)
		}
	}

	log.Println("registro da aposta inserido com sucesso", betsLoad)

	return nil
}
