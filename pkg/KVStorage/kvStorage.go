package KVStorage

import c "context"

const (
	METHODGET    = "GET"
	METHODDELETE = "DELETE"
	METHODPUT    = "PUT"
)

type Message struct {
	Ctx    c.Context   `json:"ctx"`
	Method string      `json:"method"`
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
}

type Storage struct {
	store map[string]interface{}
}

func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]interface{}),
	}
}

func (s *Storage) Get(ctx c.Context, k string) (interface{}, error) {
	return s.store[k], nil
}

func (s *Storage) Put(ctx c.Context, k string, v interface{}) error {
	s.store[k] = v
	return nil
}

func (s *Storage) Delete(ctx c.Context, k string) error {
	delete(s.store, k)
	return nil
}

func (s Storage) Store() map[string]interface{} {
	return s.store
}
