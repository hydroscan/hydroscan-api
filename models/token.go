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
	TotalSupply    string          `gorm:"column:total_supply" json:"totalSupply"`
	HoldersCount   uint64          `gorm:"column:holders_count" json:"holdersCount"`
	PriceUSD       decimal.Decimal `gorm:"column:price_usd;type:decimal(32,18)" json:"priceUSD"`
	PriceUpdatedAt time.Time       `gorm:"column:price_updated_at" json:"priceUpdatedAt"`

	Volume24h       decimal.Decimal `gorm:"column:volume_24h;type:decimal(32,18)" json:"volume24h"`
	Volume7d        decimal.Decimal `gorm:"column:volume_7d;type:decimal(32,18)" json:"volume7d"`
	Volume1m        decimal.Decimal `gorm:"column:volume_1m;type:decimal(32,18)" json:"volume1m"`
	VolumeAll       decimal.Decimal `gorm:"column:volume_all;type:decimal(32,18)" json:"volumeAll"`
	Volume24hChange float32         `gorm:"column:volume_24h_change" json:"volume24hChange"`
	Volume7dChange  float32         `gorm:"column:volume_7d_change" json:"volume7dChange"`
	Volume1mChange  float32         `gorm:"column:volume_1m_change" json:"volume1mChange"`

	// GORM ignore
	Trades24h        uint64          `gorm:"-" json:"trades24h"`
	Trades24hChange  float32         `gorm:"-" json:"trades24hChange"`
	Traders24h       uint64          `gorm:"-" json:"traders24h"`
	Traders24hChange float32         `gorm:"-" json:"traders24hChange"`
	Amount24h        decimal.Decimal `gorm:"-" json:"amount24h"`
}

func (Token) TableName() string {
	return "tokens"
}
