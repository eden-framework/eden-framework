package api

import "gitee.com/eden-framework/strings"

type Api struct {
	ServiceName string                    `json:"name"`
	Operators   map[string]*OperatorGroup `json:"operators"`
	Models      map[string]*OperatorModel `json:"models"`
	Enums       map[string]Enum           `json:"enums"`
}

func NewApi() Api {
	return Api{
		Operators: make(map[string]*OperatorGroup),
		Models:    make(map[string]*OperatorModel),
		Enums:     make(map[string]Enum),
	}
}

func (a *Api) AddGroup(name string) *OperatorGroup {
	if _, ok := a.Operators[name]; !ok {
		group := NewOperatorGroup(name, str.ToLowerSlashCase(name))
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

func (a *Api) AddEnum(id string, e Enum) {
	if _, ok := a.Enums[id]; !ok {
		a.Enums[id] = e
	}
}
