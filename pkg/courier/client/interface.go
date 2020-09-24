package client

import (
	"github.com/eden-framework/eden-framework/pkg/courier"
)

type IRequest interface {
	Do() courier.Result
}
