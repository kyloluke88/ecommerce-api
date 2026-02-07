package productskupackage

import (
	"api/app/models"
	"api/pkg/database"
)

type ProductSkuPackage struct {
	models.BaseModel
	Volume float64 `json:"volume"`
	Title  string  `json:"title"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight int     `json:"weight"`
	SkuID  int64   `json:"skuId"` // 外键
	models.CommonTimestampsField
}

func (productSkuPackageModel *ProductSkuPackage) Create() {
	database.DB.Create(&productSkuPackageModel)
}

func BatchCreatePackages(skus []ProductSkuPackage) error {
	return database.DB.CreateInBatches(skus, 100).Error
}
