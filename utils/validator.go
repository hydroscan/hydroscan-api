package utils

import (
	"regexp"
	"strings"

	"github.com/hydroscan/hydroscan-api/models"
)

func IsAddress(str string) (ret bool) {
	matched, err := regexp.MatchString("^0x[A-F0-9a-f]{40}$", str)
	if err != nil {
		ret = false
		return
	}

	ret = matched
	return
}

func IsTransaction(str string) (ret bool) {
	matched, err := regexp.MatchString("^0x[0-9a-f]{64}$", str)
	if err != nil {
		ret = false
		return
	}

	ret = matched
	return
}

func IsToken(address string) (bool, string) {
	address = strings.ToLower(address)
	token := models.Token{}
	if err := models.DB.Where("LOWER(address) = ?", address).First(&token).Error; err != nil {
		return false, ""
	}
	return true, token.Address
}

func IsTrader(address string) (bool, string) {
	address = strings.ToLower(address)
	trade := models.Trade{}
	if err := models.DB.Where("LOWER(maker_address) = ? OR LOWER(taker_address) = ?", address, address).First(&trade).Error; err != nil {
		return false, ""
	}
	if strings.ToLower(trade.MakerAddress) == address {
		return true, trade.MakerAddress
	}
	return true, trade.TakerAddress
}

func IsRelayer(address string) (bool, string) {
	address = strings.ToLower(address)
	relayer := models.Relayer{}
	if err := models.DB.Where("LOWER(address) = ?", address).First(&relayer).Error; err != nil {
		return false, ""
	}
	return true, relayer.Address
}

func IsTokenName(str string) (bool, string) {
	str = strings.ToLower(str)
	token := models.Token{}
	if err := models.DB.Where("LOWER(name) = ?", str).First(&token).Error; err != nil {
		return false, ""
	}
	return true, token.Address
}

func IsTokenSymbol(str string) (bool, string) {
	str = strings.ToLower(str)
	token := models.Token{}
	if err := models.DB.Where("LOWER(symbol) = ?", str).First(&token).Error; err != nil {
		return false, ""
	}
	return true, token.Address
}
