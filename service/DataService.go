package service

import (
	"sensor_iot/Util"
	"sensor_iot/domain"
)

func AddData(data *domain.DataModel) bool {
	Util.MyDataBase.Create(data)
	return true
}

func FindData(count int) []domain.DataModel {
	return make([]domain.DataModel, 0)
}
