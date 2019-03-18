package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/hydroscan/hydroscan-api/task"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func GetTokens(c *gin.Context) {
	query := TokensQuery{1, 25, "24H", "", "", ""}
	c.BindQuery(&query)
	log.Info("query ", query)

	page := query.Page
	pageSize := query.PageSize
	offset := (page - 1) * pageSize

	var res TokensRes
	if query.Keyword != "" {
		res = getTokensByKeyword(page, pageSize, offset, query.Keyword)
	} else if query.RelayerAddress != "" {
		res = getTokensByRelayer(page, pageSize, offset, query.RelayerAddress)
	} else if query.TraderAddress != "" {
		res = getTokensByTrader(page, pageSize, offset, query.TraderAddress)
	} else {
		res = getTokensDefault(page, pageSize, offset, query.Filter)
	}

	c.JSON(200, res)

}

func GetToken(c *gin.Context) {
	address := c.Params.ByName("address")
	token := models.Token{}
	if err := models.DB.Where("address = ?", address).First(&token).Error; gorm.IsRecordNotFoundError(err) {
		c.JSON(404, token)
	} else {
		tradesData := task.GetTrades24hData(token.Address)
		token.Trades24h = tradesData.Trades24h
		token.Trades24hChange = tradesData.Trades24hChange
		token.Traders24h = tradesData.Traders24h
		token.Traders24hChange = tradesData.Traders24hChange
		token.Amount24h = tradesData.Amount24h

		c.JSON(200, token)
	}
}
