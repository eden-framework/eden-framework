package api

type Api struct {
	Operators map[string]*OperatorGroup `json:"operators"`
	Models    map[string]*OperatorModel `json:"models"`
}

func NewApi() Api {
	return Api{
		Operators: make(map[string]*OperatorGroup),
		Models:    make(map[string]*OperatorModel),
	}
}

func (a *Api) AddGroup(name string) *OperatorGroup {
	if _, ok := a.Operators[name]; !ok {
		group := NewOperatorGroup(name)
		a.Operators[group.Name] = group
	}

	return a.Operators[name]
}

func (a *Api) GetGroup(name string) *OperatorGroup {
	return a.Operators[name]
}

func (a *Api) AddModel(name string) *OperatorModel {
	if _, ok := a.Models[name]; !ok {
		model := NewOperatorModel(name)
		a.Models[model.Name] = &model
	}

	return a.Models[name]
}
