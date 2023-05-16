package models

type Journals struct {
	UserPermission struct {
		UserName string `json:"userName"`
		IsRead   bool   `json:"isRead"`
		IsChange bool   `json:"isChange"`
	} `json:"userPermission"`
	DataJournals []struct {
		ID              int     `json:"id"`
		DT              string  `json:"dt"`
		DtAmkr          string  `json:"dtAmkr"`
		NBf             string  `json:"nBf"`
		NBrigade        string  `json:"nBrigade"`
		NShift          string  `json:"nShift"`
		ProdDaily       int     `json:"prodDaily"`
		ProdShift       float64 `json:"prodShift"`
		SiAVGShift      float64 `json:"siAVGShift"`
		SiAVGDaily      float64 `json:"siAVGDaily"`
		StateShift      bool    `json:"stateShift"`
		IDData          int     `json:"idData"`
		StMasterSmeni   string  `json:"stMasterSmeni"`
		SmMasterPechi   string  `json:"smMasterPechi"`
		Gazovschik      string  `json:"gazovschik"`
		MashinistShihty string  `json:"mashinistShihty"`
		Dispatcher      string  `json:"dispatcher"`
		InfoAdditional1 string  `json:"infoAdditional1"`
		InfoAdditional2 string  `json:"infoAdditional2"`
	} `json:"dataJournals"`
}
