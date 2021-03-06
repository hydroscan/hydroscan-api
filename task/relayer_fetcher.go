package task

import (
	"io/ioutil"
	"os"

	"github.com/hydroscan/hydroscan-api/internal/json"
	"github.com/hydroscan/hydroscan-api/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type RelayerInfo struct {
	Url     string `json:"url"`
	Address string `json:"address"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
}

func UpdateRelayers() {
	log.Info("UpdateRelayers")

	relayers := getRelayers()
	log.Info(relayers)

	for _, r := range relayers {
		mRelayer := models.Relayer{}
		if err := models.DB.Where("address = ?", r.Address).First(&mRelayer).Error; gorm.IsRecordNotFoundError(err) {
			mRelayer = models.Relayer{
				Address: r.Address,
				Url:     r.Url,
				Name:    r.Name,
				Slug:    r.Slug,
			}
			models.DB.Create(&mRelayer)
		} else {
			models.DB.Model(&mRelayer).Updates(models.Relayer{
				Url:  r.Url,
				Name: r.Name,
				Slug: r.Slug,
			})
		}
	}
}

func CreateRelayerByAddressIfNotExist(address string) {
	mRelayer := models.Relayer{}
	if err := models.DB.Where("address = ?", address).First(&mRelayer).Error; gorm.IsRecordNotFoundError(err) {
		mRelayer = models.Relayer{
			Address: address,
			Url:     "https://etherscan.io/address/" + address,
			Name:    "Relayer-" + address[0:6] + "..." + address[len(address)-4:], // like Relayer-0x4949...33FD
			Slug:    address,
		}
		models.DB.Create(&mRelayer)
	}
}

func getRelayers() []RelayerInfo {
	jsonFile, err := os.Open(ResourcePath + "/relayers.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	relayers := []RelayerInfo{}

	json.Unmarshal(byteValue, &relayers)
	return relayers
}
