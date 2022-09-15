package api

type uploadType int

const (
	newApp uploadType = iota
	appUpdate
)
