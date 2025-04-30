package state

import (
	"atlas-reactors/data/reactor/item"
)

type Model struct {
	theType      int32
	reactorItem  *item.Model
	activeSkills []uint32
	nextState    int8
}
