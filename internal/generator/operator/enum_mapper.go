package operator

import (
	"gitee.com/eden-framework/enumeration"
	str "gitee.com/eden-framework/strings"
	"strings"
)

var serviceEnumMap = map[string]map[string]enumeration.Enum{}

func RegisterEnum(serviceName string, tpe string, options ...enumeration.EnumOption) {
	serviceName = strings.ToLower(str.ToUpperCamelCase(serviceName))
	if serviceEnumMap[serviceName] == nil {
		serviceEnumMap[serviceName] = map[string]enumeration.Enum{}
	}
	serviceEnumMap[serviceName][tpe] = options
}

func GetEnumByServiceName(serviceName string) map[string]enumeration.Enum {
	serviceName = strings.ToLower(str.ToUpperCamelCase(serviceName))
	if serviceEnumMap[serviceName] == nil {
		serviceEnumMap[serviceName] = map[string]enumeration.Enum{}
	}
	return serviceEnumMap[serviceName]
}
