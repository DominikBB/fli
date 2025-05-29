package repository_test

import (
	"os"
	"testing"

	"github.com/dominikbb/fli/internal/filestore"
	"github.com/dominikbb/fli/internal/repository"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type Store interface {
	repository.Creator
	repository.Getter
}

func TestStoreImplementations(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		setup     func(store Store)
		testFunc  func(store Store) error
		expectErr error
	}{
		{
			name: "Store with duplicate key",
			setup: func(store Store) {
				_ = store.Store("value1", "key1")
			},
			testFunc: func(store Store) error {
				return store.Store("value2", "key1")
			},
			expectErr: repository.ErrDuplicate,
		},
		{
			name:  "Get with non-existent key",
			setup: func(store Store) {},
			testFunc: func(store Store) error {
				_, err := store.Get("nonexistent")
				return err
			},
			expectErr: repository.ErrNotFound,
		},
	}

	// Define implementations
	implementations := []struct {
		name  string
		store Store
	}{
		{
			name:  "MemoryStore",
			store: repository.NewMemoryStore(),
		},
		{
			name:  "FileStore",
			store: setupFileStore(t),
		},
	}

	// Run tests for all implementations
	for _, impl := range implementations {
		for _, tt := range tests {
			t.Run(impl.name+"_"+tt.name, func(t *testing.T) {
				tt.setup(impl.store)
				err := tt.testFunc(impl.store)
				assert.Equal(t, tt.expectErr, err)
			})
		}
	}
}

func setupFileStore(t *testing.T) *filestore.Store {
	// Create a temporary fli.json file
	fileName := "fli.json"
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("failed to create fli.json: %v", err)
	}
	_, err = file.WriteString("{}")
	if err != nil {
		t.Fatalf("failed to write valid JSON to fli.json: %v", err)
	}
	file.Close()

	// Set up Viper
	v := viper.New()
	v.SetConfigFile(fileName)
	v.SetConfigType("json")

	// Clean up the file after tests
	t.Cleanup(func() {
		os.Remove(fileName)
	})

	// Create and return the FileStore
	store, err := filestore.NewStore("")
	if err != nil {
		t.Fatalf("failed to create FileStore: %v", err)
	}
	return store
}
