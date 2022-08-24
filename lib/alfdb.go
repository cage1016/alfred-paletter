package lib

import (
	"context"
)

type ClipBoard struct {
	DataHash string `gorm:"column:dataHash"`
	Item     string `gorm:"column:item"`
	AppPath  string `gorm:"column:apppath"`
	DateTime string `gorm:"column:ts"`
}

func (ClipBoard) TableName() string {
	return "clipboard"
}

type ClipBoardRepository interface {
	List(ctx context.Context, limit int) ([]*ClipBoard, error)
}
