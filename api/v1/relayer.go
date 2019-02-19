package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
)

func GetRelayers(c *gin.Context) {
	var relayers []models.Relayer
	if err := models.DB.Find(&relayers).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, relayers)
	}
}

func GetRelayer(c *gin.Context) {
	slug := c.Params.ByName("slug")
	relayer := models.Relayer{}
	if err := models.DB.Where("slug = ?", slug).First(&relayer).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, relayer)
	}
}
