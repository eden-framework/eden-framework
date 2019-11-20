package api

type Api struct {
	Operators map[string]*OperatorGroup `json:"operators"`
	Models    string                    `json:"models"`
}

func NewApi() Api {
	return Api{
		Operators: make(map[string]*OperatorGroup),
		Models:    "",
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
