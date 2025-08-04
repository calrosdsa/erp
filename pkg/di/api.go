package di

import (
	"context"
	"fmt"
)

func Get(ctx context.Context, key string) any {
	ctn, ok := ctx.Value(containerKey).(*container)
	fmt.Println("CONTAIONER",ctn)
	if !ok {
		panic("container does not exist on context")
	}

	return ctn.Get(key)
}
