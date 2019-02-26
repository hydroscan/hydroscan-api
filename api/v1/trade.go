package apiv1

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/redis"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func GetTrades(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	i, err := strconv.ParseInt(pageQuery, 10, 64)
	page := int(i)
	if err != nil {
		page = 1
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var trades []models.Trade
	if err := models.DB.Order("block_number desc").Order("log_index desc").Offset(offset).Limit(pageSize).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").Find(&trades).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		type resType struct {
			Page      int            `json:"page"`
			TotalPage int            `json:"totalPage"`
			Trades    []models.Trade `json:"trades"`
		}
		res := resType{page, 0, trades}
		totalCount := 0
		models.DB.Table("trades").Count(&totalCount)
		res.TotalPage = int(math.Ceil(float64(totalCount) / float64(pageSize)))

		c.JSON(200, res)
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
	if err := models.DB.Where("uuid = ?", uuid).Preload("Relayer").Preload("BaseToken").Preload("QuoteToken").First(&trade).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, trade)
	}
}

func GetTradesChart(c *gin.Context) {
	filter := c.DefaultQuery("filter", "1M")
	var res []struct {
		Dt    time.Time       `json:"date"`
		Sum   decimal.Decimal `json:"volume"`
		Count int             `json:"trades"`
	}
	trunc := "day"
	from := time.Now().Add(-30 * 24 * time.Hour)
	switch filter {
	case "24H":
		trunc = "hour"
		from = time.Now().Add(-24 * time.Hour)
	case "7D":
		trunc = "hour"
		from = time.Now().Add(-7 * 24 * time.Hour)
	case "1M":
		trunc = "day"
		from = time.Now().Add(-30 * 24 * time.Hour)
	case "1Y":
		trunc = "day"
		from = time.Now().Add(-365 * 24 * time.Hour)
	case "ALL":
		trunc = "day"
		from = time.Now().Add(-1000 * 24 * time.Hour)
	default:
		c.AbortWithStatus(404)
		return
	}
	models.DB.Raw("select date_trunc(?, date) as dt, sum(volume_usd), count(1) from trades where date >= ? group by dt order by dt", trunc, from).Scan(&res)
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
