package courier

import (
	"context"
	"fmt"
)

func GetContextValue(ctx context.Context, contextProvider IContextProvider) interface{} {
	v := ctx.Value(contextProvider.ContextKey())
	if v == nil {
		panic(fmt.Errorf("context providor %#v required", contextProvider))
	}
	return v
}
