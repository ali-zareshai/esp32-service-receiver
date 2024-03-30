package domain

import (
	"gorm.io/gorm"
	"sensor_iot/Util"
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
	Util.MyDataBase.Create(model)
}
