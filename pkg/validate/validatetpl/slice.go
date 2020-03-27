package validatetpl

import (
	"fmt"
	"reflect"
)

func NewValidateSlice(min uint64, max uint64, elemValidateFunc func(v interface{}) (bool, string)) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			sliceValue := reflect.ValueOf(v)
			elemCount := sliceValue.Len()
			if uint64(elemCount) < min || (max > 0 && uint64(elemCount) > max) {
				return false, fmt.Sprintf(SLICE_ELEM_NUM_NOT_IN_RANGE, min, max, elemCount)
			}
			if elemValidateFunc != nil {
				for i := 0; i < elemCount; i++ {
					elemValue := sliceValue.Index(i)
					if ok, errMsg := elemValidateFunc(elemValue.Interface()); !ok {
						return false, fmt.Sprintf(SLICE_ELEM_INVALID, errMsg)
					}
				}
			}
			return true, ""
		}
		return false, TYPE_NOT_SLICE
	}
}
