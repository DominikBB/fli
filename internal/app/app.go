package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dominikbb/fli/internal/filestore"
	"github.com/dominikbb/fli/internal/repository"
	"github.com/spf13/viper"
)

type App struct {
	Creator repository.Creator
	Getter  repository.Getter
}

func NewApp() *App {
	globalStoreLocation := viper.GetString("globalStoreLocation")
	if globalStoreLocation == "" {
		home, _ := os.UserHomeDir()
		globalStoreLocation = filepath.Join(home, ".cache", "fli.json")
	}

	fileStore, err := filestore.NewStore(globalStoreLocation)
	if err != nil {
		panic(fmt.Errorf("failed to create file store: %w", err))
	}

	return &App{
		Creator: fileStore,
		Getter:  fileStore,
	}
}

func NewMemoryApp() *App {
	store := repository.NewMemoryStore()
	return &App{
		Creator: store,
		Getter:  store,
	}
}
