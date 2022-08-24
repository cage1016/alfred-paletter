package sqlite

import (
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Home string `envconfig:"HOME"`
	Path string `envconfig:"QS_DB_PATH" default:"Library/Application Support/Alfred/Databases"`
	Name string `envconfig:"QS_DB_NAME" default:"clipboard.alfdb"`
}

func Connect(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(cfg.Home, cfg.Path, cfg.Name)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
