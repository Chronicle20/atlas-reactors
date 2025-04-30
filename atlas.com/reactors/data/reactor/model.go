package reactor

import (
	"atlas-reactors/data/reactor/point"
	"atlas-reactors/data/reactor/state"
)

type Model struct {
	tl          point.Model
	br          point.Model
	stateInfo   map[int8][]state.Model
	timeoutInfo map[int8]int32
}
