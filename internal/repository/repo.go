package repository

import (
	"errors"
	"strings"
)

type (
	Creator interface {
		Store(string, ...string) error
	}

	Getter interface {
		Get(...string) (string, error)
		List(...string) ([][]string, error)
	}
)

var (
	ErrDuplicate = errors.New("another command is associated with these tags")
	ErrNotFound  = errors.New("no command found for that set of tags")
)

var (
	_ Creator = (*MemoryStore)(nil)
	_ Getter  = (*MemoryStore)(nil)
)

const KeySeparator = " "

type MemoryStore struct {
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStore) Store(value string, keys ...string) error {
	if len(keys) == 0 {
		return errors.New("at least one key is required")
	}
	compositeKey := strings.Join(keys, KeySeparator)
	if _, exists := m.data[compositeKey]; exists {
		return ErrDuplicate
	}
	m.data[compositeKey] = value
	return nil
}

func (m *MemoryStore) Get(keys ...string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("at least one key is required")
	}
	compositeKey := strings.Join(keys, KeySeparator)
	value, exists := m.data[compositeKey]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}

func (m *MemoryStore) List(keys ...string) ([][]string, error) {
	substring := strings.Join(keys, KeySeparator)
	var results [][]string
	for key, value := range m.data {
		if strings.Contains(key, substring) {
			results = append(results, []string{key, value})
		}
	}
	return results, nil
}
