package sqlite

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/cage1016/alfred-paletter/lib"
)

type clipBoardRepository struct {
	mu sync.RWMutex
	db *gorm.DB
}

func New(db *gorm.DB) lib.ClipBoardRepository {
	return &clipBoardRepository{
		mu: sync.RWMutex{},
		db: db,
	}
}

func (repo *clipBoardRepository) List(ctx context.Context, limit int) (res []*lib.ClipBoard, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	// err = repo.db.WithContext(ctx).Where("dataType", "1").Order("ts desc").Limit(limit).Find(&res).Error
	err = repo.db.Raw("SELECT dataHash,item,apppath,strftime('%Y-%m-%d %H:%M',ts+978307200,'unixepoch','localtime') as ts from clipboard WHERE dataType = 1 ORDER BY ts DESC LIMIT ?", limit).Scan(&res).Error
	return
}
