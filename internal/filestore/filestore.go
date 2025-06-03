package filestore

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dominikbb/fli/internal/repository"
	"github.com/spf13/viper"
)

const KeySeparator = " "

type Store struct {
	viper *viper.Viper
}

var (
	_ repository.Creator = (*Store)(nil)
	_ repository.Getter  = (*Store)(nil)
	_ repository.Remover = (*Store)(nil)
)

func NewStore(filePath string) (*Store, error) {
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("json")

	if err := v.ReadInConfig(); err != nil {
		// Try to handle multiple kinds of "not found" errors
		var notFound viper.ConfigFileNotFoundError
		var pathErr *fs.PathError

		if errors.As(err, &notFound) || errors.As(err, &pathErr) {
			// Treat both as missing file cases
			dir := filepath.Dir(filePath)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return nil, fmt.Errorf("failed to create parent directory: %w", err)
				}
			}
			if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
				return nil, fmt.Errorf("failed to write initial json: %w", err)
			}
			if err := v.ReadInConfig(); err != nil {
				return nil, fmt.Errorf("failed to read config after creation: %w", err)
			}
		} else {
			// Real error — config is present but corrupted or unreadable
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	return &Store{viper: v}, nil
}

func (f *Store) Store(value string, keys ...string) error {
	if len(keys) == 0 {
		return errors.New("at least one key is required")
	}
	compositeKey := strings.Join(keys, KeySeparator)
	if f.viper.IsSet(compositeKey) && f.viper.GetString(compositeKey) != "" {
		return repository.ErrDuplicate
	}
	f.viper.Set(compositeKey, value)
	return f.viper.WriteConfig()
}

func (f *Store) Get(keys ...string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("at least one key is required")
	}
	compositeKey := strings.Join(keys, KeySeparator)
	value := f.viper.GetString(compositeKey)
	if value == "" {
		return "", repository.ErrNotFound
	}
	return value, nil
}

func (f *Store) List(keys ...string) ([][]string, error) {
	substring := strings.Join(keys, KeySeparator)
	var results [][]string
	for _, key := range f.viper.AllKeys() {
		if f.viper.GetString(key) == "" {
			continue
		}
		if strings.Contains(key, substring) {
			results = append(results, []string{key, f.viper.GetString(key)})
		}
	}
	return results, nil
}

func (f *Store) Delete(keys ...string) error {
	if len(keys) == 0 {
		return errors.New("at least one key is required")
	}
	compositeKey := strings.Join(keys, KeySeparator)
	if !f.viper.IsSet(compositeKey) {
		return repository.ErrNotFound
	}
	f.viper.Set(compositeKey, nil)
	f.viper.Set(compositeKey, "")
	return f.viper.WriteConfig()
}
