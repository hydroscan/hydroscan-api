package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Trade struct {
	gorm.Model
	UUID               string          `gorm:"column:uuid;unique_index;not null" json:"uuid"`
	BlockNumber        uint64          `gorm:"column:block_number;index;unique_index:idx_block_number_log_index" json:"blockNumber"`
	BlockHash          string          `gorm:"column:block_hash" json:"blockHash"`
	TransactionHash    string          `gorm:"column:transaction_hash" json:"transactionHash"`
	LogIndex           uint            `gorm:"column:log_index;;unique_index:idx_block_number_log_index" json:"logIndex"`
	Date               time.Time       `gorm:"column:date;index" json:"date"`
	QuoteTokenPriceUSD decimal.Decimal `gorm:"column:quote_token_price_usd;type:decimal(32,18)" json:"quoteTokenPriceUSD"`
	VolumeUSD          decimal.Decimal `gorm:"column:volume_usd;type:decimal(32,18)" json:"volumeUSD"`
	BaseTokenAddress   string          `gorm:"column:base_token_address;index" json:"baseTokenAddress"`
	QuoteTokenAddress  string          `gorm:"column:quote_token_address;index" json:"quoteTokenAddress"`
	RelayerAddress     string          `gorm:"column:relayer_address;index" json:"relayerAddress"`
	MakerAddress       string          `gorm:"column:maker_address" json:"makerAddress"`
	TakerAddress       string          `gorm:"column:taker_address" json:"takerAddress"`
	BaseTokenAmount    decimal.Decimal `gorm:"column:base_token_amount;type:decimal(32,18)" json:"baseTokenAmount"`
	QuoteTokenAmount   decimal.Decimal `gorm:"column:quote_token_amount;type:decimal(32,18)" json:"quoteTokenAmount"`
	MakerFee           decimal.Decimal `gorm:"column:maker_fee;type:decimal(32,18)" json:"makerFee"`
	TakerFee           decimal.Decimal `gorm:"column:taker_fee;type:decimal(32,18)" json:"takerFee"`
	MakerGasFee        decimal.Decimal `gorm:"column:maker_gas_fee;type:decimal(32,18)" json:"makerGasFee"`
	MakerRebate        decimal.Decimal `gorm:"column:maker_rebate;type:decimal(32,18)" json:"makerRebate"`
	TakerGasFee        decimal.Decimal `gorm:"column:taker_gas_fee;type:decimal(32,18)" json:"takerGasFee"`
	BaseToken          Token           `gorm:"foreignkey:base_token_address;association_foreignkey:address" json:"baseToken"`
	QuoteToken         Token           `gorm:"foreignkey:quote_token_address;association_foreignkey:address" json:"quoteToken"`
	Relayer            Relayer         `gorm:"foreignkey:relayer_address;association_foreignkey:address" json:"relayer"`
}

func (Trade) TableName() string {
	return "trades"
}

func (m *Trade) BeforeCreate(scope *gorm.Scope) (err error) {
	m.UUID = fmt.Sprint(uuid.New())
	return
}
