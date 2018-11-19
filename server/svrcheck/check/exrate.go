package check

type Exrate interface {
	Check

	GetPrice(baseCur string, targetCur string) (float64, error)
}
