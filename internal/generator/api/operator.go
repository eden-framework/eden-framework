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

type OperatorMethod struct {
	Group     *OperatorGroup `json:"-"`
	Name      string         `json:"name"`
	inputDef  map[string]bool
	outputDef map[string]bool
	Inputs    map[string]OperatorModel `json:"inputs"`
	Outputs   map[string]OperatorModel `json:"outputs"`
}

func NewOperatorMethod(group *OperatorGroup, name string) *OperatorMethod {
	return &OperatorMethod{
		Group:     group,
		Name:      name,
		inputDef:  make(map[string]bool),
		outputDef: make(map[string]bool),
		Inputs:    make(map[string]OperatorModel),
		Outputs:   make(map[string]OperatorModel),
	}
}

func (m *OperatorMethod) AddInputDef(name string) {
	m.inputDef[name] = true
}

func (m *OperatorMethod) AddOutputDef(name string) {
	m.outputDef[name] = true
}

func (m *OperatorMethod) AddInput(model OperatorModel) {
	m.Inputs[model.Name] = model
}

func (m *OperatorMethod) AddOutput(model OperatorModel) {
	m.Outputs[model.Name] = model
}
