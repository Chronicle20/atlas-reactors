package reactor

import (
	"atlas-reactors/data/reactor/point"
	state2 "atlas-reactors/data/reactor/state"
	"github.com/Chronicle20/atlas-model/model"
	"strconv"
)

type RestModel struct {
	Id          uint32                      `json:"-"`
	TL          point.RestModel             `json:"tl"`
	BR          point.RestModel             `json:"br"`
	StateInfo   map[int8][]state2.RestModel `json:"stateInfo"`
	TimeoutInfo map[int8]int32              `json:"timeoutInfo"`
}

func (r RestModel) GetName() string {
	return "reactors"
}

func (r RestModel) GetID() string {
	return strconv.Itoa(int(r.Id))
}

func (r *RestModel) SetID(strId string) error {
	id, err := strconv.Atoi(strId)
	if err != nil {
		return err
	}
	r.Id = uint32(id)
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
	si := make(map[int8][]state2.Model)
	for k, vs := range rm.StateInfo {
		si[k] = make([]state2.Model, 0)
		for _, v := range vs {
			sm, err := state2.Extract(v)
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
