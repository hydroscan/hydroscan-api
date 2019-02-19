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
	Trades24h uint            `gorm:"column:trades_24h" json:"trades24h"`
	Volume24h decimal.Decimal `gorm:"column:volume_24h;type:decimal(32,18)" json:"volume24h"`
}

func (Relayer) TableName() string {
	return "relayers"
}
