package transform

type ParameterMap map[string]*ParameterMeta

func (m ParameterMap) List() (list []*ParameterMeta) {
	for _, parameterMeta := range m {
		list = append(list, parameterMeta)
	}
	return
}

func (m ParameterMap) Add(parameterMeta *ParameterMeta) {
	m[parameterMeta.Field.Name] = parameterMeta
}

func (m ParameterMap) Get(fieldName string) (rv *ParameterMeta, ok bool) {
	rv, ok = m[fieldName]
	return
}

func (m ParameterMap) Len() int {
	return len(m)
}
