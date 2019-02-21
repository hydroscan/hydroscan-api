package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Token struct {
	gorm.Model
	Name            string          `json:"name"`
	Symbol          string          `json:"symbol"`
	Decimals        uint            `json:"decimals"`
	Address         string          `gorm:"unique_index" json:"address"`
	PriceUSD        decimal.Decimal `gorm:"column:price_usd;type:decimal(32,18)" json:"priceUSD"`
	PriceUpdatedAt  time.Time       `json:"priceUpdatedAt"`
	Volume24h       decimal.Decimal `gorm:"column:volume_24h;type:decimal(32,18)" json:"volume24h"`
	Volume7d        decimal.Decimal `gorm:"column:volume_7d;type:decimal(32,18)" json:"volume7d"`
	VolumeAll       decimal.Decimal `gorm:"column:volume_all;type:decimal(32,18)" json:"volumeAll"`
	Volume24hChange float32         `gorm:"column:volume_24h_change" json:"volume24hChange"`
	Volume7dChange  float32         `gorm:"column:volume_7d_change" json:"volume7dChange"`
}

func (Token) TableName() string {
	return "tokens"
}
