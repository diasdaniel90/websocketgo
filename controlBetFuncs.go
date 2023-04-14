package main

import (
	"log"
	"time"
)

const maxGale = 2

func sinal2Playbet(sliceSignals *[]msgSignalStruct,
	msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct,
) {
	log.Println("Executando a função após 4 segundos...", sliceSignals, msgStatusRec)

	if len(*sliceSignals) != 0 {
		for _, value := range *sliceSignals {
			bet := betBotStruct{
				idBet:     msgStatusRec.idBet,
				timestamp: msgStatusRec.timestamp,
				color:     value.Color,
				source:    value.Source,
				win:       false,
				status:    false,
				gale:      0,
			}
			*sliceBets = append(*sliceBets, bet)

			log.Println("value", value)
		}

		*sliceSignals = []msgSignalStruct{}
	}

	log.Println("apostas feitas", sliceBets)
}

func validateBet(msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct) {
	if msgStatusRec.betStatus != waiting && len(*sliceBets) != 0 {
		log.Print("valida win")

		sliceBetsGale := []betBotStruct{}

		for index := range *sliceBets {
			if (*sliceBets)[index].idBet == msgStatusRec.idBet && (*sliceBets)[index].color == msgStatusRec.color {
				(*sliceBets)[index].win = true
				log.Println("vaivai", (*sliceBets)[index].win)
			} else if (*sliceBets)[index].gale < maxGale {
				sliceBetsGale = append(sliceBetsGale, (*sliceBets)[index])
				sliceBetsGale[len(sliceBetsGale)-1].gale++
				log.Println("loss vai no gale", sliceBetsGale)
			}

			err := saveToDatabaseBets(&(*sliceBets)[index])
			if err != nil {
				log.Printf("error no banco: %s", err)

				continue
			}
		}

		log.Println("Color:", msgStatusRec.color)
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

func controlBet(msgStatusChan <-chan msgStatusStruct, msgSignalChan <-chan msgSignalStruct) {
	log.Println("###########11##################")

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
				setID(&sliceBets, msgStatusRec)

				time.AfterFunc(tempoEspera*time.Second, func() {
					go sinal2Playbet(&sliceSignals, msgStatusRec, &sliceBets)
				})
			}

			validateBet(msgStatusRec, &sliceBets)

		case signalMsg, ok := <-msgSignalChan:
			if !ok {
				log.Println("Canal msgSignalChan fechado")

				continue
			}
			// log.Println("Recebeu sinal msgSignalChan", mensagens)
			// log.Println("Recebeu sinal ", signalMsg)
			sliceSignals = append(sliceSignals, signalMsg)

		default:
			// Faça algo aqui se ambos os canais estiverem vazios.
			// Por exemplo, tente novamente mais tarde.
			time.Sleep(100 * time.Millisecond)
		}
	}
}
