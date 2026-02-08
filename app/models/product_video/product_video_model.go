package productvideo

import "api/app/models"

type ProductVideo struct {
	models.BaseModel
	ProductID uint64 `json:"product_id" gorm:"uniqueIndex"`

	VideoID  int64 `gorm:"uniqueIndex"`
	Title    string
	CoverURL string
	VideoURL string
	State    int16

	models.CommonTimestampsField
}
