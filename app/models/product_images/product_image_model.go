package productimages

import (
	"api/app/models"
	"api/pkg/database"
)

type ProductImage struct {
	models.BaseModel
	ProductID           uint64
	FullPathImageURI    string `json:"fullPathImageURI"`
	SearchImageURI      string `json:"searchImageURI"`
	Size220x220ImageURI string `json:"size220x220ImageURI" gorm:"column:size_220x220_image_uri"`
	Size310x310ImageURI string `json:"size310x310ImageURI" gorm:"column:size_310x310_image_uri"`
	SummImageURI        string `json:"summImageURI"`

	models.CommonTimestampsField
}

func (imageModel *ProductImage) Create() {
	database.DB.Create(&imageModel)
}

func BatchCreateImages(skus []ProductImage) error {
	return database.DB.CreateInBatches(skus, 100).Error
}
