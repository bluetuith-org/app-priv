package views

type viewID int

const (
	viewAdTree viewID = 1 << iota
	viewTabs
	viewInfo
	viewOperations
)
