package task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenPrice struct {
	Rate decimal.Decimal `json:"rate"`
	TS   int64           `json:"ts"`
}

type TokenInfo struct {
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Decimals          string
	DecimalsInterface interface{} `json:"decimals"` // ethplorer api return decimals type can be string or number
	TotalSupply       string      `json:"totalSupply"`
	HoldersCount      uint64      `json:"holdersCount"`
	Price             TokenPrice  `json:"price"`
}

func GetToken(address string) models.Token {
	mToken := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&mToken).Error; gorm.IsRecordNotFoundError(err) {
		tokenInfo := GetTokenInfo(address)

		decimals, err := strconv.ParseUint(tokenInfo.Decimals, 10, 64)
		if err != nil {
			panic(err)
		}

		mToken = models.Token{
			Address:        address,
			Decimals:       uint(decimals),
			Name:           tokenInfo.Name,
			Symbol:         tokenInfo.Symbol,
			TotalSupply:    tokenInfo.TotalSupply,
			HoldersCount:   tokenInfo.HoldersCount,
			PriceUSD:       tokenInfo.Price.Rate,
			PriceUpdatedAt: time.Unix(tokenInfo.Price.TS, 0),
		}

		models.DB.Create(&mToken)
	}
	return mToken
}

func UpdateTokenPrices() {
	log.Info("UpdateTokenPrices")

	mTokens := []models.Token{}
	models.DB.Find(&mTokens)

	for _, mToken := range mTokens {
		tokenInfo := GetTokenInfo(mToken.Address)

		models.DB.Model(&mToken).Updates(models.Token{
			Name:           tokenInfo.Name,
			Symbol:         tokenInfo.Symbol,
			TotalSupply:    tokenInfo.TotalSupply,
			HoldersCount:   tokenInfo.HoldersCount,
			PriceUSD:       tokenInfo.Price.Rate,
			PriceUpdatedAt: time.Unix(tokenInfo.Price.TS, 0),
		})
	}
}

func GetTokenInfo(address string) TokenInfo {
	// old TUSD, ethplorer doesn't have this token
	if address == "0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E" {
		// TUSD
		address = "0x0000000000085d4780B73119b644AE5ecd22b376"
	}

	// https://changelog.makerdao.com/releases/mainnet/1.0.0/contracts.json
	// new DAI, MCD_DAI, ethplorer doesn't have this token
	if address == "0x6B175474E89094C44Da98b954EedeAC495271d0F" {
		// old DAI
		address = "0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359"
	}

	tokenInfo := TokenInfo{}

	if address == "0x000000000000000000000000000000000000000E" {
		tokenInfo.Symbol = "ETH"
		tokenInfo.Name = "Ether"
		tokenInfo.Decimals = "18"
	} else {
		log.Info("GetTokenInfo ", address)
		url := "http://api.ethplorer.io/getTokenInfo/" + address + "?apiKey=" + viper.GetString("ethplorer_apikey")
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal([]byte(body), &tokenInfo)
	}

	// ethplorer api return decimals type can be string or number
	if tokenInfo.Decimals == "" {
		tokenInfo.Decimals = fmt.Sprint(tokenInfo.DecimalsInterface)
	}
	// get ETH price for WETH and ETH
	if tokenInfo.Symbol == "WETH" || tokenInfo.Symbol == "ETH" {
		lastPrice := getETHLastPrice()
		log.Info("WETH Price: ", tokenInfo.Price)
		log.Info("ETH Price: ", lastPrice)
		if !lastPrice.Rate.IsZero() {
			tokenInfo.Price = lastPrice
		}
	}

	// old DAI => SAI
	if address == "0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359" {
		tokenInfo.Name = "SAI"
		tokenInfo.Symbol = "SAI"
	}

	// these tokens' name and symbol are not string
	// https://etherscan.io/token/0x9469d013805bffb7d3debe5e7839237e535ec483#readContract
	if address == "0x9469D013805bFfB7D3DEBe5E7839237e535ec483" {
		tokenInfo.Name = "Evolution Land Global Token"
		tokenInfo.Symbol = "RING"
	}
	// https://etherscan.io/token/0xeb269732ab75A6fD61Ea60b06fE994cD32a83549#readContract
	if address == "0xeb269732ab75A6fD61Ea60b06fE994cD32a83549" {
		tokenInfo.Name = "USDx"
		tokenInfo.Symbol = "USDx"
	}
	// https://etherscan.io/token/0x431ad2ff6a9C365805eBaD47Ee021148d6f7DBe0#readContract
	if address == "0x431ad2ff6a9C365805eBaD47Ee021148d6f7DBe0" {
		tokenInfo.Name = "dForce"
		tokenInfo.Symbol = "DF"
	}
	// https://etherscan.io/token/0x2630997aAB62fA1030a8b975e1AA2dC573b18a13#readContract
	if address == "0x2630997aAB62fA1030a8b975e1AA2dC573b18a13" {
		tokenInfo.Name = "HYPE Token"
		tokenInfo.Symbol = "HYPE"
	}

	return tokenInfo
}

func getETHLastPrice() TokenPrice {
	url := "http://api.ethplorer.io/getTop?apiKey=" + viper.GetString("ethplorer_apikey")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type Result struct {
		Tokens []struct {
			Address string     `json:"address"`
			Name    string     `json:"name"`
			Symbol  string     `json:"symbol"`
			Price   TokenPrice `json:"price"`
		} `json:"tokens"`
	}

	result := Result{}
	json.Unmarshal([]byte(body), &result)
	for _, token := range result.Tokens {
		if token.Symbol == "ETH" {
			return token.Price
		}
	}

	return TokenPrice{}
}
