package domain

import (
	"gorm.io/gorm"
)

type DataModel struct {
	gorm.Model
	Device string  `gorm:"column:device" json:"device"`
	Result float64 `gorm:"column:result" json:"result"`
	M32    string  `gorm:"column:m32" json:"m32"`
	M33    string  `gorm:"column:m33" json:"m33"`
	Ref    string  `gorm:"column:ref" json:"ref"`
	Vout   string  `gorm:"column:vout" json:"vout"`
}
