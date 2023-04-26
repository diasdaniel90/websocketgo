package main

type betBotStruct struct {
	idBet          string
	timestamp      int64
	timestampSinal int64
	color          int
	source         string
	win            bool
	status         string
	gale           int
	amount         float32
	balanceWin     float32
}

type msgStatusStruct struct {
	idBet     string
	timestamp int64
	betStatus string
	color     int
	betRoll   int
}

type msgSignalStruct struct {
	Type      string  `json:"type"`
	Timestamp float64 `json:"time"`
	Color     int     `json:"color"`
	Source    int     `json:"source"`
}
type lastMsgStruct struct {
	lastUpdatedAt string
	lastID        string
	lastIDWaiting string
}
type payloadStruct struct {
	IDBet                string            `json:"id"`
	Color                int               `json:"color"`
	Roll                 int               `json:"roll"`
	CreatedAt            string            `json:"created_at"`
	Timestamp            int64             `json:"timestamp"`
	UpdatedAt            string            `json:"updated_at"`
	Status               string            `json:"status"`
	TotalRedEurBet       float64           `json:"total_red_eur_bet"`
	TotalRedBetsPlaced   int               `json:"total_red_bets_placed"`
	TotalWhiteEurBet     float64           `json:"total_white_eur_bet"`
	TotalWhiteBetsPlaced int               `json:"total_white_bets_placed"`
	TotalBlackEurBet     float64           `json:"total_black_eur_bet"`
	TotalBlackBetsPlaced int               `json:"total_black_bets_placed"`
	TotalBetsPlaced      int               `json:"totalBetsPlaced"`
	TotalEurBet          float64           `json:"totalEurBet"`
	TotalRetentionEur    float64           `json:"totalRetentionEur"`
	Bets                 []betsUsersStruct `json:"bets"`
}

type betsUsersStruct struct {
	IDBetUser    string  `json:"id"`
	Color        int     `json:"color"`
	Amount       float32 `json:"amount"`
	CurrencyType string  `json:"currency_type"`
	Status       string  `json:"status"`
	User         struct {
		IDStr string `json:"id_str"`
	} `json:"user"`
}

type configStruct struct {
	EnvRef        string `json:"envRef"`
	MySQLDatabase string `json:"mySqlDatabase"`
	MySQLUser     string `json:"mySqlUser"`
	MySQLPassword string `json:"mySqlPassword"`
	MySQLHost     string `json:"mySqlHost"`
	MySQLPort     string `json:"mySqlPort"`
}
