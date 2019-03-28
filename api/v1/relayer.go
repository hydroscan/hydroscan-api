package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/hydroscan/hydroscan-api/utils"
	"github.com/jinzhu/gorm"
)

func GetRelayers(c *gin.Context) {
	var relayers []models.Relayer
	if err := models.DB.Order("volume_24h DESC").Find(&relayers).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, relayers)
	} else {

		for i, _ := range relayers {
			tradesData := task.GetRelayerTrades24hData(relayers[i].Address)
			relayers[i].Traders24h = tradesData.Traders24h
			relayers[i].Trades24h = tradesData.Trades24h
		}
		c.JSON(200, relayers)
	}
}

func GetRelayer(c *gin.Context) {
	slug := c.Params.ByName("slug")
	relayer := models.Relayer{}

	statment := models.DB.Where("slug = ?", slug)
	if utils.IsAddress(slug) {
		if isTrue, key := utils.IsRelayer(slug); isTrue {
			statment = models.DB.Where("address = ?", key)
		} else {
			c.JSON(404, relayer)
		}
	}

	if err := statment.First(&relayer).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, relayer)
	} else {
		res := getRelayerService(relayer)
		c.JSON(200, res)
	}
}
