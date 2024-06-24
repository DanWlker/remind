package entity

type FlagEntity struct {
	Name      string
	Shorthand string
	Usage     string
}

type BoolFlagEntity struct {
	FlagEntity
	Value bool
}
