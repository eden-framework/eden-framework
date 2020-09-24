package transform

import (
	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
)

type ParameterErrors struct {
	StatusError *status_error.StatusError
}

func (p *ParameterErrors) Merge(err error) {
	if err == nil {
		return
	}
	statusError := status_error.FromError(err)
	if p.StatusError == nil {
		p.StatusError = statusError
		return
	}
	p.StatusError = p.StatusError.WithErrorFields(statusError.ErrorFields...)
}

func (p *ParameterErrors) Err() error {
	if p.StatusError != nil {
		return p.StatusError
	}
	return nil
}
