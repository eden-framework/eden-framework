package api

import "strings"

type OperatorModel struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Package   string          `json:"package"`
	Fields    []OperatorField `json:"fields,omitempty"`
	NeedAlias bool            `json:"needAlias"`
}

func NewOperatorModel(name string, pkgID string) OperatorModel {
	return OperatorModel{
		ID:      strings.Join([]string{pkgID, name}, "."),
		Name:    name,
		Package: pkgID,
		Fields:  make([]OperatorField, 0),
	}
}

func (m *OperatorModel) AddField(key, keyType, tag, alias, ipt string, pointer bool) {
	m.Fields = append(m.Fields, OperatorField{
		Key:     key,
		Type:    keyType,
		Tag:     tag,
		Alias:   alias,
		Imports: ipt,
		Pointer: pointer,
	})
}

func (m *OperatorModel) WalkFields(walker func(f OperatorField)) {
	for _, field := range m.Fields {
		walker(field)
	}
}

type OperatorField struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Tag     string `json:"tag"`
	Alias   string `json:"alias"`
	Imports string `json:"imports"`
	Pointer bool   `json:"pointer"`
}
