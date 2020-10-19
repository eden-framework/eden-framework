package project

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "github.com/eden-framework/enumeration"
)

var InvalidShippingProcessorType = errors.New("invalid ShippingProcessorType")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("ShippingProcessorType", map[string]string{
		"ALIYUN_REGISTRY": "aliyun registry",
	})
}

func ParseShippingProcessorTypeFromString(s string) (ShippingProcessorType, error) {
	switch s {
	case "":
		return SHIPPING_PROCESSOR_TYPE_UNKNOWN, nil
	case "ALIYUN_REGISTRY":
		return SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY, nil
	}
	return SHIPPING_PROCESSOR_TYPE_UNKNOWN, InvalidShippingProcessorType
}

func ParseShippingProcessorTypeFromLabelString(s string) (ShippingProcessorType, error) {
	switch s {
	case "":
		return SHIPPING_PROCESSOR_TYPE_UNKNOWN, nil
	case "aliyun registry":
		return SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY, nil
	}
	return SHIPPING_PROCESSOR_TYPE_UNKNOWN, InvalidShippingProcessorType
}

func (ShippingProcessorType) EnumType() string {
	return "ShippingProcessorType"
}

func (ShippingProcessorType) Enums() map[int][]string {
	return map[int][]string{
		int(SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY): {"ALIYUN_REGISTRY", "aliyun registry"},
	}
}

func (v ShippingProcessorType) String() string {
	switch v {
	case SHIPPING_PROCESSOR_TYPE_UNKNOWN:
		return ""
	case SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY:
		return "ALIYUN_REGISTRY"
	}
	return "UNKNOWN"
}

func (v ShippingProcessorType) Label() string {
	switch v {
	case SHIPPING_PROCESSOR_TYPE_UNKNOWN:
		return ""
	case SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY:
		return "aliyun registry"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*ShippingProcessorType)(nil)

func (v ShippingProcessorType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidShippingProcessorType
	}
	return []byte(str), nil
}

func (v *ShippingProcessorType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseShippingProcessorTypeFromString(string(bytes.ToUpper(data)))
	return
}
