package api

type Api struct {
	ServiceName string                    `json:"name"`
	Operators   map[string]*OperatorGroup `json:"operators"`
	Models      map[string]*OperatorModel `json:"models"`
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

func (a *Api) WalkOperators(walker func(g *OperatorGroup)) {
	for _, group := range a.Operators {
		walker(group)
	}
}

func (a *Api) AddModel(model *OperatorModel) {
	if _, ok := a.Models[model.ID]; !ok {
		a.Models[model.ID] = model
	}
}
