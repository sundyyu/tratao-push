package check

type Exrate interface {
	Loop

	Update() error
	GetPrice(baseCur string, targetCur string) (float64, error)
}
