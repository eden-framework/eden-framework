package api

func NewTimeModel() OperatorModel {
	return OperatorModel{
		ID:        "time.Time",
		Name:      "Time",
		Package:   "time",
		NeedAlias: true,
	}
}
