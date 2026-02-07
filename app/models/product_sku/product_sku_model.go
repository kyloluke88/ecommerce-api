package productsku

import (
	"api/app/models"
	"api/pkg/database"

	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

type ProductSku struct {
	models.BaseModel
	Title         string          `json:"titme"`
	Description   string          `json:"description"`
	ProductId     uint64          `json:"product_id"`
	Stock         int64           `json:"stock"`
	SkuID         int64           `json:"sku_id" gorm:"uniqueIndex"`
	SaleCount     int64           `json:"sale_count"`
	PromotionSku  bool            `json:"promotion_sku"`
	Price         decimal.Decimal `json:"price"`
	DiscountPrice decimal.Decimal `json:"discount_price"`

	models.CommonTimestampsField
}

func (productSku *ProductSku) Create() {
	database.DB.Create(&productSku)
}

func BatchCreateSKUs(skus []ProductSku) error {

	return database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sku_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"price",
			"discount_price",
			"stock",
		}),
	}).CreateInBatches(skus, 100).Error
}
