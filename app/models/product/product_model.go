package product

import (
	"api/app/models"
	"api/pkg/database"

	"github.com/shopspring/decimal"
)

type Product struct {
	models.BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string
	MaxPrice    decimal.Decimal
	MinPrice    decimal.Decimal
	OnSale      bool
	CompanyID   uint64
	SaleCount   int64
	models.CommonTimestampsField
}

func (productModel *Product) Create() {
	database.DB.Create(&productModel)
}
