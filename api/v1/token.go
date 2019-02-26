package apiv1

import (
	"math"
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
		page = 1
	}
	pageSize := 25
	offset := (page - 1) * pageSize

	var tokens []models.Token
	if err := models.DB.Offset(offset).Limit(pageSize).Find(&tokens).Error; gorm.IsRecordNotFoundError(err) {
		c.AbortWithStatus(404)
	} else {
		type resType struct {
			Page      int            `json:"page"`
			PageSize  int            `json:"pageSize"`
			TotalPage int            `json:"totalPage"`
			Count     uint64         `json:"count"`
			Tokens    []models.Token `json:"tokens"`
		}
		res := resType{page, pageSize, 0, 0, tokens}
		models.DB.Table("tokens").Count(&res.Count)
		res.TotalPage = int(math.Ceil(float64(res.Count) / float64(pageSize)))

		c.JSON(200, res)
	}
}

func GetTokensTop(c *gin.Context) {
	filter := c.DefaultQuery("filter", "24H")

	order := "volume_24h desc"
	switch filter {
	case "24H":
		order = "volume_24h desc"
	case "7D":
		order = "volume_7d desc"
	case "ALL":
		order = "volume_all desc"
	default:
		c.AbortWithStatus(404)
		return
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
