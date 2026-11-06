package item

type Type int

const (
	TypeEnergyCell Type = iota
	TypeMeds
	TypeFood
	TypeCloth
	TypeBook
	TypeToy
	TypeMisc
	TypeOther // last
)
