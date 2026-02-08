package product

import (
	"api/app/models"
	productattr "api/app/models/product_attr"
	productimages "api/app/models/product_images"
	productsku "api/app/models/product_sku"
	productvideo "api/app/models/product_video"
	"api/pkg/database"

	"github.com/shopspring/decimal"
)

type Product struct {
	models.BaseModel
	Title         string `json:"title"`
	Description   string `json:"description"`
	Location      string
	MaxPrice      decimal.Decimal
	MinPrice      decimal.Decimal
	OnSale        bool
	CompanyID     uint64
	SaleCount     int64
	ProductId1688 int64 `json:"product_id_1688" gorm:"column:product_id_1688"`

	// nil = 没视频（语义清晰）
	// Preload("Video") 时，不存在不会报错
	// JSON 序列化时可自然忽略
	// 如果你 不打算创建视频 不要给 Video 赋 &ProductVideo{}（因为会插空记录）
	Video  *productvideo.ProductVideo
	Images []productimages.ProductImage
	Skus   []productsku.ProductSku
	Attrs  []productattr.ProductAttr

	models.CommonTimestampsField
}

func (productModel *Product) Create() {
	database.DB.Create(&productModel)
}
