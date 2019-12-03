package api

import (
	"github.com/profzone/eden-framework/pkg/enumeration"
	"sort"
)

type Enum enumeration.Enum

func (enum Enum) Sort() Enum {
	sort.Slice(enum, func(i, j int) bool {
		switch enum[i].Value.(type) {
		case string:
			return enum[i].Value.(string) < enum[j].Value.(string)
		case int64:
			return enum[i].Value.(int64) < enum[j].Value.(int64)
		case float64:
			return enum[i].Value.(float64) < enum[j].Value.(float64)
		}
		return true
	})
	return enum
}

func (enum Enum) Labels() (labels []string) {
	for _, e := range enum {
		labels = append(labels, e.Label)
	}
	return
}

func (enum Enum) Vals() (vals []interface{}) {
	for _, e := range enum {
		vals = append(vals, e.Val)
	}
	return
}

func (enum Enum) Values() (values []interface{}) {
	for _, e := range enum {
		values = append(values, e.Value)
	}
	return
}
