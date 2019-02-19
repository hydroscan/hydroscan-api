package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Token struct {
	gorm.Model
	Name           string          `json:"name"`
	Symbol         string          `json:"symbol"`
	Decimals       uint            `json:"decimals"`
	Address        string          `gorm:"unique_index" json:"address"`
	PriceUSD       decimal.Decimal `gorm:"column:price_usd;type:decimal(32,18)" json:"priceUSD"`
	PriceUpdatedAt time.Time       `json:"priceUpdatedAt"`
	Trades24h      uint            `gorm:"column:trades_24h" json:"trades24h"`
	Volume24h      decimal.Decimal `gorm:"column:volume_24h;type:decimal(32,18)" json:"volume24h"`
}

func (Token) TableName() string {
	return "tokens"
}
