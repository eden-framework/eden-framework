package api

type OperatorGroup struct {
	Name    string                     `json:"name"`
	Path    string                     `json:"path"`
	Methods map[string]*OperatorMethod `json:"methods"`
	IsPush  bool                       `json:"isPush"`
}

func NewOperatorGroup(name, path string) *OperatorGroup {
	return &OperatorGroup{
		Name:    name,
		Path:    path,
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
	Path    string         `json:"path"`
	Inputs  []string       `json:"inputs"`
	Outputs []string       `json:"outputs"`
}

func NewOperatorMethod(group *OperatorGroup, name, path string) *OperatorMethod {
	return &OperatorMethod{
		Group:   group,
		Name:    name,
		Path:    path,
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
