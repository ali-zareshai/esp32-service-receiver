package domain

import (
	"encoding/json"
	"gorm.io/gorm"
	"os"
	"sensor_iot/utils"
	"strconv"
)

type DataModel struct {
	gorm.Model
	Device string  `gorm:"column:device" json:"device"`
	Result float64 `gorm:"column:result" json:"result"`
	M32    float64 `gorm:"column:m32" json:"m32"`
	M33    float64 `gorm:"column:m33" json:"m33"`
	Ref    float64 `gorm:"column:ref" json:"ref"`
	Vout   float64 `gorm:"column:vout" json:"vout"`
}

func (model *DataModel) Save() {
	utils.MyDataBase.Create(model)
}

func (model *DataModel) ToJson() string {
	str, err := json.Marshal(model)
	if err != nil {
		return ""
	}
	return string(str)
}

func (model *DataModel) AfterSave(tx *gorm.DB) error {
	alert, err := strconv.ParseFloat(os.Getenv("ALERT_RESULT"), 64)
	if err != nil {
		alert = 15.0
	}
	if model.Result >= alert {
		jsonInfo := model.ToJson()
		utils.PublishRedis(utils.AlertRedisChannels, jsonInfo)
	}
	return nil
}
