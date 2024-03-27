package service

import "sensor_iot/domain"

func AddData(data domain.DataJsonRequest) bool {

	return true
}

func FindData(count int) []domain.DataJsonRequest {
	return make([]domain.DataJsonRequest, 0)
}
