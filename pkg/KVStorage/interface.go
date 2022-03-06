package KVStorage

import "context"

type KVStorage interface {
	Get(context.Context, string) (interface{}, error)
	Put(context.Context, string, interface{}) error
	Delete(context.Context, string) error
}
