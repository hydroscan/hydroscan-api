package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/hydroscan/hydroscan-api/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func GetTrades(c *gin.Context) {
	query := TradesQuery{1, 25, "", "", "", "", "", ""}
	c.BindQuery(&query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var trades []models.Trade
	statment := models.DB.Table("trades").Order("block_number desc").Order("log_index desc")
	if query.Transaction != "" {
		statment = statment.Where("transaction_hash = ?", query.Transaction)
	} else if query.BaseTokenAddress != "" && query.QuoteTokenAddress != "" {
		statment = statment.Where("base_token_address = ? AND quote_token_address = ?", query.BaseTokenAddress, query.QuoteTokenAddress)
	} else if query.TokenAddress != "" {
		statment = statment.Where("base_token_address = ? OR quote_token_address = ?", query.TokenAddress, query.TokenAddress)
	} else if query.TraderAddress != "" {
		statment = statment.Where("maker_address = ? OR taker_address = ?", query.TraderAddress, query.TraderAddress)
	} else if query.RelayerAddress != "" {
		statment = statment.Where("relayer_address = ?", query.RelayerAddress)
	}

	type resType struct {
		Page     int            `json:"page"`
		PageSize int            `json:"pageSize"`
		Count    uint64         `json:"count"`
		Trades   []models.Trade `json:"trades"`
	}
	if err := statment.Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		res := resType{page, pageSize, 0, trades}
		c.JSON(404, res)
	} else {
		res := resType{page, pageSize, 0, trades}
		statment.Count(&res.Count)

		c.JSON(200, res)
	}
}

func GetTrade(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	trade := models.Trade{}
	if err := models.DB.Where("uuid = ?", uuid).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").First(&trade).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, trade)
	} else {
		c.JSON(200, trade)
	}
}

func GetTrader(c *gin.Context) {
	address := c.Params.ByName("address")
	trade := models.Trade{}
	if err := models.DB.Where("maker_address = ? OR taker_address = ?", address, address).First(&trade).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, TraderRes{})
		return
	}

	res := getTraderService(address)
	c.JSON(200, res)
}

func GetTradesChart(c *gin.Context) {
	query := TradesChartQuery{"1M", "", "", ""}
	c.BindQuery(&query)

	log.Info(query)
	res := getTradesChartService(query)
	c.JSON(200, res)
}

func GetTradesIndicators(c *gin.Context) {
	res, err := redis.Client.Get("indicators").Result()
	if err != nil {
		panic(err)
	}
	indicators := task.Indicators{}
	json.Unmarshal([]byte(res), &indicators)
	c.JSON(200, indicators)
}

func GetTradesSearch(c *gin.Context) {
	keyword := c.Query("keyword")

	var res struct {
		SearchType string `json:"searchType"`
		SearchKey  string `json:"searchKey"`
	}

	if utils.IsAddress(keyword) {
		if isTrue, searchKey := utils.IsToken(keyword); isTrue {
			res.SearchType = "TOKEN"
			res.SearchKey = searchKey

		} else if isTrue, searchKey := utils.IsRelayer(keyword); isTrue {
			res.SearchType = "RELAYER"
			res.SearchKey = searchKey

		} else if isTrue, searchKey := utils.IsTrader(keyword); isTrue {
			res.SearchType = "TRADER"
			res.SearchKey = searchKey

		}
	} else if utils.IsTransaction(keyword) {
		res.SearchType = "TRANSACTION"
		res.SearchKey = keyword

	} else {
		token := models.Token{}
		if err := models.DB.Where("name ILIKE ? or symbol ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").Order("volume_24h desc").First(&token).Error; err == nil {
			res.SearchType = "TOKENS"
			res.SearchKey = keyword
		}
	}

	c.JSON(200, res)
}
