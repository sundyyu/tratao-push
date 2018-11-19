package check

type Check interface {
	Update()
	Loop()
}
