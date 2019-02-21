package models

import (
	"github.com/jinzhu/gorm"
)

type Relayer struct {
	gorm.Model
	Name    string `json:"name"`
	Url     string `json:"url"`
	Slug    string `gorm:"unique_index" json:"slug"`
	Address string `gorm:"unique_index" json:"address"`
}

func (Relayer) TableName() string {
	return "relayers"
}
