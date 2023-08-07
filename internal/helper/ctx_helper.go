package helper

import (
	"context"
	"fmt"
	"os"

	"github.com/quick-im/quick-im-core/internal/contant"
)

type ctxOpt struct {
	key contant.ContentKey
	val any
}

func CtxOptWarp[T any](key contant.ContentKey, val T) ctxOpt {
	return ctxOpt{
		key: key,
		val: val,
	}
}

func InitCtx(parentCtx context.Context, opts ...ctxOpt) context.Context {
	for i := range opts {
		parentCtx = context.WithValue(parentCtx, opts[i].key, opts[i].val)
	}
	return parentCtx
}

func GetCtxValue[T any](ctx context.Context, key contant.ContentKey, assertType T) T {
	v := ctx.Value(key)
	if v == nil {
		fmt.Fprintf(os.Stderr, "The service requires a dependency on the context \"%s\"\n", key)
		os.Exit(1)
	}
	_, ok := v.(T)
	if !ok {
		fmt.Fprintf(os.Stderr, "Assertion type error \"%T\"\n", assertType)
		os.Exit(1)
	}
	return v.(T)
}
