package project

//go:generate eden generate enum --type-name=ShippingProcessorType
// api:enum
type ShippingProcessorType uint8

// 打包处理器类型
const (
	SHIPPING_PROCESSOR_TYPE_UNKNOWN          ShippingProcessorType = iota
	SHIPPING_PROCESSOR_TYPE__ALIYUN_REGISTRY                       // aliyun registry
)
