package models

type Tapping struct {
	ID               int     `json:"id"`
	IDJournal        int     `json:"idJournal"`
	NumTaphole       int     `json:"numTaphole"`
	NumTapping       int     `json:"numTapping"`
	DtOpenTaphole    string  `json:"dtOpenTaphole"`
	DtCloseTaphole   string  `json:"dtCloseTaphole"`
	EnumLadleTapping string  `json:"enumLadleTapping"`
	ListLaldes       []Ladle `json:"listLaldes"`
	TapholeLength    float64 `json:"tapholeLength"`
	TapholeState     string  `json:"tapholeState"`
	MarkBlow         string  `json:"markBlow"`
	Charge           string  `json:"charge"`
	WeightSum        float64 `json:"weightSum"`
	Temper           float64 `json:"temper"`
	CountLaldleTap   int     `json:"countLaldleTap"`
	Proba            float64 `json:"proba"`
	DtChem           string  `json:"dtChem"`
	Si               float64 `json:"si"`
	Mn               float64 `json:"mn"`
	S                float64 `json:"s"`
	P                float64 `json:"p"`
	C                float64 `json:"c"`
	LadleDestAnother int     `json:"ladleDestAnother"`
	LadleDestBof     int     `json:"ladleDestBof"`
	LadleDestMc      int     `json:"ladleDestMc"`
	LadleDestPm      int     `json:"ladleDestPm"`
	LadleDestCpi     int     `json:"ladleDestCpi"`
}

type Ladle struct {
	ID          int      `json:"id"`
	IDTapping   int      `json:"idTapping"`
	Ladle       string   `json:"ladle"`
	Weight      float64  `json:"weight"`
	IsDelivered bool     `json:"isDelivered"`
	Destination string   `json:"destination"`
	DtInsert    string   `json:"dtInsert"`
	Chemical    Chemical `json:"chemical"`
}

type Chemical struct {
	ID         int     `json:"id"`
	IDJournal  int     `json:"idJournal"`
	IDTapping  int     `json:"idTapping"`
	NTapping   int     `json:"nTapping"`
	NBf        int     `json:"nBf"`
	Ladle      string  `json:"ladle"`
	Proba      float64 `json:"proba"`
	NumTaphole int     `json:"numTaphole"`
	Dt         string  `json:"dt"`
	Si         float64 `json:"si"`
	Mn         float64 `json:"mn"`
	S          float64 `json:"s"`
	P          float64 `json:"p"`
	Cr         float64 `json:"cr"`
	Ti         float64 `json:"ti"`
	C          float64 `json:"c"`
	T          float64 `json:"t"`
	SiO2       float64 `json:"siO2"`
	Al2O3      float64 `json:"al2O3"`
	CaO        float64 `json:"caO"`
	MnO        float64 `json:"mnO"`
	MgO        float64 `json:"mgO"`
	FeO        float64 `json:"feO"`
	Ss         float64 `json:"ss"`
	Och1       float64 `json:"och1"`
	Och2       float64 `json:"och2"`
	FeAll      float64 `json:"feAll"`
	K2o        float64 `json:"k2o"`
	Na2O       float64 `json:"na2O"`
	ZnO        float64 `json:"znO"`
	TiO2       float64 `json:"tiO2"`
	DtInsert   string  `json:"dtInsert"`
	DtUpdate   string  `json:"dtUpdate"`
	CreatedOn  string  `json:"createdOn"`
	ModifyBy   string  `json:"modifyBy"`
}
