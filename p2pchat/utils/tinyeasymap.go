package utils

import "fmt"

type TinyEasyMap struct {
	data map[string]any
}

func NewEasyMap(m map[string]any) TinyEasyMap {
	if m == nil {
		panic("map must not be nil")
	}
	return TinyEasyMap{m}
}

func (m TinyEasyMap) Get(key string) (any, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not exists", key)
	}
	return val, nil
}

func (m TinyEasyMap) GetString(key string) (string, error) {
	val, err := m.Get(key)
	if err != nil {
		return "", err
	}
	strval, ok := val.(string)
	if ok {
		return strval, nil
	}
	return "", fmt.Errorf("invalid type key='%s'", key)
}

func (m TinyEasyMap) GetUint64(key string) (uint64, error) {
	val, err := m.Get(key)
	if err != nil {
		return 0, err
	}
	return ConvUint64(val)
}
