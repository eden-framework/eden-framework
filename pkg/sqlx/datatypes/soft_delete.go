package datatypes

import (
	"github.com/eden-framework/eden-framework/pkg/enumeration"
)

type SoftDelete struct {
	Enabled enumeration.Bool `db:"F_enabled" sql:"int(8) unsigned NOT NULL DEFAULT '1'" json:"-"`
}

func (e *SoftDelete) Enable() {
	e.Enabled = enumeration.BOOL__TRUE
}

func (e *SoftDelete) Disable() {
	e.Enabled = enumeration.BOOL__FALSE
}
