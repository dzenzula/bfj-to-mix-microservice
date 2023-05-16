package models

type MixJournals struct {
	UserPermission struct {
		UserName string `json:"userName"`
		IsRead   bool   `json:"isRead"`
		IsChange bool   `json:"isChange"`
	} `json:"userPermission"`
	DataJournals []struct {
		ID                   int     `json:"id"`
		DT                   string  `json:"dt"`
		DtAmkr               string  `json:"dtAmkr"`
		NMix                 string  `json:"nMix"`
		NBrigade             string  `json:"nBrigade"`
		NShift               string  `json:"nShift"`
		ProdDaily            int     `json:"prodDaily"`
		ProdShift            float64 `json:"prodShift"`
		SiAVGShift           float64 `json:"siAVGShift"`
		SiAVGDaily           float64 `json:"siAVGDaily"`
		StateShift           bool    `json:"stateShift"`
		IDData               int     `json:"idData"`
		MasterProizvUchastka string  `json:"masterProizvUchastka"`
		Brigadir             string  `json:"brigadir"`
		Mikserovoj           string  `json:"mikserovoj"`
		InfoAdditional1      string  `json:"infoAdditional1"`
		InfoAdditional2      string  `json:"infoAdditional2"`
	} `json:"dataJournals"`
}
