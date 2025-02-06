package data

import (
	"atlas-reactors/reactor/data/point"
	"atlas-reactors/reactor/data/state"
	"github.com/Chronicle20/atlas-model/model"
)

type RestModel struct {
	Id          string                     `json:"-"`
	TL          point.RestModel            `json:"tl"`
	BR          point.RestModel            `json:"br"`
	StateInfo   map[int8][]state.RestModel `json:"stateInfo"`
	TimeoutInfo map[int8]int32             `json:"timeoutInfo"`
}

func (r RestModel) GetName() string {
	return "reactors"
}

func (r RestModel) GetID() string {
	return r.Id
}

func (r *RestModel) SetID(id string) error {
	r.Id = id
	return nil
}

func Extract(rm RestModel) (Model, error) {
	tl, err := model.Map(point.Extract)(model.FixedProvider(rm.TL))()
	if err != nil {
		return Model{}, err
	}
	br, err := model.Map(point.Extract)(model.FixedProvider(rm.BR))()
	if err != nil {
		return Model{}, err
	}
	si := make(map[int8][]state.Model)
	for k, vs := range rm.StateInfo {
		si[k] = make([]state.Model, 0)
		for _, v := range vs {
			sm, err := state.Extract(v)
			if err != nil {
				return Model{}, err
			}

			si[k] = append(si[k], sm)
		}
	}

	return Model{
		tl:          tl,
		br:          br,
		stateInfo:   si,
		timeoutInfo: rm.TimeoutInfo,
	}, nil
}
