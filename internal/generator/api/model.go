package api

type OperatorModel struct {
	Name    string          `json:"name"`
	Imports [][]string      `json:"imports"`
	Fields  []OperatorField `json:"fields"`
}

func NewOperatorModel(name string) OperatorModel {
	return OperatorModel{
		Name:    name,
		Imports: make([][]string, 0),
		Fields:  make([]OperatorField, 0),
	}
}

func (m *OperatorModel) AddField(key, keyType string) {
	m.Fields = append(m.Fields, OperatorField{
		Key:  key,
		Type: keyType,
	})
}

func (m *OperatorModel) AddImport(path, alias string) {
	m.Imports = append(m.Imports, []string{alias, path})
}

type OperatorField struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
