package model

const (
	PtListPrice = "LP" // LIST PRICE
	PtCashBuy   = "CB" // CASH BUY
	PtCashSell  = "CS" // CASH SELL
	PtBuy       = "B"  // BUY
	PtSell      = "S"  // SELL
)

func IsValidPriceType(pt string) bool {
	return pt == PtListPrice || pt == PtCashBuy || pt == PtCashSell ||
		pt == PtBuy || pt == PtSell
}

type Alarm struct {
	Id            int64   `json:"id"`
	Account       string  `json:"account"`
	BaseCur       string  `json:"basecur"`
	TargetCur     string  `json:"targetcur"`
	Price         float64 `json:"price"`
	LbPrice       float64 `json:"lbprice"`
	UbPrice       float64 `json:"ubprice"`
	Enabled       bool    `json:"enabled"`
	DeviceId      string  `json:"devid"`
	DeviceOS      string  `json:"devos"`
	DeviceCountry string  `json:"devcountry"`
	DeviceLang    string  `json:"devlang"`
	AppKey        string  `json:"appkey"`
	Ltt           int64   `json:"ltt"` // last trigger time
	UpdateTime    int64   `json:"updatetime"`
	CreateTime    int64   `json:"createtime"`
}

func GetFields() []string {
	fields := []string{
		"Id",
		"Account",
		"BaseCur",
		"TargetCur",
		"Price",
		"LbPrice",
		"UbPrice",
		"Enabled",
		"DeviceId",
		"DeviceOS",
		"DeviceCountry",
		"DeviceLang",
		"AppKey",
		"Ltt",
		"UpdateTime",
		"CreateTime"}

	return fields
}
