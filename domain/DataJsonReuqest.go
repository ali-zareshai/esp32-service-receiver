package domain

type DataJsonRequest struct {
	Device string  `json:"device" binding:"required"`
	Result float64 `json:"result" binding:"required"`
}
