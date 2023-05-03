package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	maxGale     = 4
	amount      = 2.0
	tempoEspera = 4
)

func sinal2Playbet(sliceSignals *[]msgSignalStruct,
	msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct,
) {
	log.Printf("Executando a função após %d segundos...", tempoEspera)

	log.Println("sliceSignals", sliceSignals)
	log.Println("msgStatusRec", msgStatusRec)

	if len(*sliceSignals) != 0 {
		for _, value := range *sliceSignals {
			bet := betBotStruct{
				idBet:          msgStatusRec.idBet,
				timestamp:      msgStatusRec.timestamp,
				timestampSinal: int64(value.Timestamp),
				color:          value.Color,
				source:         fmt.Sprint(value.Source),
				win:            false,
				status:         "simulado",
				gale:           0,
				amount:         amount,
				balanceWin:     0.0,
			}
			*sliceBets = append(*sliceBets, bet)

			log.Println("value", value)
		}

		*sliceSignals = []msgSignalStruct{}
	}

	log.Println("apostas feitas", sliceBets)
}

func validateBet(dbConexao *sql.DB, msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct) {
	if msgStatusRec.betStatus != waiting && len(*sliceBets) != 0 {

		sliceBetsGale := []betBotStruct{}

		for index := range *sliceBets {
			if (*sliceBets)[index].idBet == msgStatusRec.idBet && (*sliceBets)[index].color == msgStatusRec.color {
				(*sliceBets)[index].win = true

				(*sliceBets)[index].balanceWin = (*sliceBets)[index].amount / amount
			} else {
				(*sliceBets)[index].balanceWin = -((*sliceBets)[index].amount / amount)
				if (*sliceBets)[index].gale < maxGale {
					sliceBetsGale = append(sliceBetsGale, (*sliceBets)[index])
					sliceBetsGale[len(sliceBetsGale)-1].gale++
					sliceBetsGale[len(sliceBetsGale)-1].amount *= 2
					log.Println("loss vai no gale", sliceBetsGale)
				}
			}
		}
		// log.Println("print", sliceBets)
		// log.Printf("O tipo de myVar é %T\n", sliceBets)
		// log.Println("print", *sliceBets)
		// log.Printf("O tipo de myVar é %T\n", *sliceBets)

		err := saveToDatabaseBets(dbConexao, sliceBets)
		if err != nil {
			log.Printf("error no banco: %s", err)
		}

		log.Println("Resultado Color:", msgStatusRec.color)
		log.Println("resultado bet", *sliceBets)
		*sliceBets = make([]betBotStruct, len(sliceBetsGale))
		copy(*sliceBets, sliceBetsGale)
	}
}

func setID(sliceBets *[]betBotStruct, msgStatusRec msgStatusStruct) {
	for index := range *sliceBets {
		(*sliceBets)[index].idBet = msgStatusRec.idBet
	}
}

func controlBet(dbConexao *sql.DB, msgStatusChan <-chan msgStatusStruct, msgSignalChan <-chan msgSignalStruct) {
	sliceSignals := []msgSignalStruct{}
	sliceBets := []betBotStruct{}
	// var valido string
	for {
		select {
		case msgStatusRec, ok := <-msgStatusChan:
			if !ok {
				log.Println("Canal msgStatusChan fechado")

				continue
			}

			if msgStatusRec.betStatus == waiting {
				log.Print("vai esperar")
				setID(&sliceBets, msgStatusRec)

				time.AfterFunc(tempoEspera*time.Second, func() {
					go sinal2Playbet(&sliceSignals, msgStatusRec, &sliceBets)
				})
			}

			validateBet(dbConexao, msgStatusRec, &sliceBets)

		case signalMsg, ok := <-msgSignalChan:
			if !ok {
				log.Println("Canal msgSignalChan fechado")

				continue
			}
			// log.Println("Recebeu sinal msgSignalChan", mensagens)
			// log.Println("Recebeu sinal ", signalMsg)
			sliceSignals = append(sliceSignals, signalMsg)

			// default:
			// 	// Faça algo aqui se ambos os canais estiverem vazios.
			// 	// Por exemplo, tente novamente mais tarde.
			// 	time.Sleep(100 * time.Millisecond)
		}
	}
}
