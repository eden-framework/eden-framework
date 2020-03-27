package client

import (
	"github.com/profzone/eden-framework/pkg/courier"
)

type IRequest interface {
	Do() courier.Result
}
