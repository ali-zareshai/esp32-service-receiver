package service

import (
	"sensor_iot/Util"
	"sensor_iot/domain"
)

func FindData(count int, device string) []domain.DataModel {
	var dataModels []domain.DataModel
	query := Util.MyDataBase.Model(&domain.DataModel{})
	if device != "" {
		query.Where("device=?", device)
	}
	query.Limit(count).Order("created_at").Find(&dataModels)
	return dataModels
}
