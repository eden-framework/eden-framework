package api

type OperatorGroup struct {
	Name    string                     `json:"name"`
	Methods map[string]*OperatorMethod `json:"methods"`
	IsPush  bool                       `json:"isPush"`
}

func NewOperatorGroup(name string) *OperatorGroup {
	return &OperatorGroup{
		Name:    name,
		Methods: make(map[string]*OperatorMethod),
	}
}

func (g *OperatorGroup) AddMethod(method *OperatorMethod) {
	method.Group = g
	g.Methods[method.Name] = method
}

func (g *OperatorGroup) AddMethods(methods ...*OperatorMethod) {
	for _, method := range methods {
		g.AddMethod(method)
	}
}

func (g *OperatorGroup) WalkMethods(walker func(m *OperatorMethod)) {
	for _, method := range g.Methods {
		walker(method)
	}
}

type OperatorMethod struct {
	Group   *OperatorGroup `json:"-"`
	Name    string         `json:"name"`
	Inputs  []string       `json:"inputs"`
	Outputs []string       `json:"outputs"`
}

func NewOperatorMethod(group *OperatorGroup, name string) *OperatorMethod {
	return &OperatorMethod{
		Group:   group,
		Name:    name,
		Inputs:  make([]string, 0),
		Outputs: make([]string, 0),
	}
}

func (m *OperatorMethod) AddInput(model *OperatorModel) {
	m.Inputs = append(m.Inputs, model.ID)
}

func (m *OperatorMethod) AddOutput(model *OperatorModel) {
	m.Outputs = append(m.Outputs, model.ID)
}

func (m *OperatorMethod) WalkInputs(walker func(i string)) {
	for _, input := range m.Inputs {
		walker(input)
	}
}

func (m *OperatorMethod) WalkOutputs(walker func(i string)) {
	for _, input := range m.Outputs {
		walker(input)
	}
}
