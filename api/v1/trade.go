package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
)

func GetTrades(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	i, err := strconv.ParseInt(pageQuery, 10, 64)
	page := int(i)
	if err != nil {
		page = 0
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var trades []models.Trade
	if err := models.DB.Order("block_number desc").Order("log_index desc").Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, trades)
	}
}

func GetTradesLatest(c *gin.Context) {
	pageSize := 8
	offset := 0
	var trades []models.Trade
	if err := models.DB.Order("block_number desc").Order("log_index desc").Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, trades)
	}
}

func GetTrade(c *gin.Context) {
	uuid := c.Params.ByName("uuid")
	trade := models.Trade{}
	if err := models.DB.Where("uuid = ?", uuid).First(&trade).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, trade)
	}
}
