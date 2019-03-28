package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Relayer struct {
	gorm.Model
	Name      string          `json:"name"`
	Url       string          `json:"url"`
	Slug      string          `gorm:"unique_index" json:"slug"`
	Address   string          `gorm:"unique_index" json:"address"`
	Volume24h decimal.Decimal `gorm:"column:volume_24h;type:decimal(32,18)" json:"volume24h"`
	// GORM ignore
	Trades24h  uint64 `gorm:"-" json:"trades24h"`
	Traders24h uint64 `gorm:"-" json:"traders24h"`
}

func (Relayer) TableName() string {
	return "relayers"
}
