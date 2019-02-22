package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
)

func GetTokens(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	i, err := strconv.ParseInt(pageQuery, 10, 64)
	page := int(i)
	if err != nil {
		page = 0
	}
	pageSize := 20
	offset := (page - 1) * pageSize

	var tokens []models.Token
	if err := models.DB.Offset(offset).Limit(pageSize).Find(&tokens).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, tokens)
	}
}

func GetTokensTop(c *gin.Context) {
	orderQuery := c.DefaultQuery("order", "volume_24h")

	order := "volume_24h desc"
	if orderQuery == "volume_24h" || orderQuery == "volume_7d" || orderQuery == "volume_all" {
		order = orderQuery + " desc"
	}

	pageSize := 10
	offset := 0

	var tokens []models.Token
	if err := models.DB.Offset(offset).Limit(pageSize).Order(order).Find(&tokens).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, tokens)
	}
}

func GetToken(c *gin.Context) {
	address := c.Params.ByName("address")
	token := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&token).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, token)
	}
}
