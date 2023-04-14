package main

import "log"

func sinal2Playbet(sliceSignals *[]msgSignalStruct,
	msgStatusRec msgStatusStruct, sliceBets *[]betBotStruct,
) {
	log.Println("Executando a função após 4 segundos...", sliceSignals, msgStatusRec)

	for _, value := range *sliceSignals {
		bet := betBotStruct{
			idBet:     msgStatusRec.idBet,
			timestamp: msgStatusRec.timestamp,
			color:     value.Color,
			source:    value.Source,
			win:       false,
			status:    false,
		}
		*sliceBets = append(*sliceBets, bet)

		log.Println("value", value)
	}

	*sliceSignals = []msgSignalStruct{}

	log.Println("apostas feitas", sliceBets)
}
