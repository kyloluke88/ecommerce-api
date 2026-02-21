package product

import (
	"api/app/models"
	productattr "api/app/models/product_attr"
	productimages "api/app/models/product_images"
	productsku "api/app/models/product_sku"
	productskupackage "api/app/models/product_sku_package"
	productvideo "api/app/models/product_video"
	"api/pkg/database"

	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	models.BaseModel
	Title         string `json:"title"`
	Description   string `json:"description"`
	SourceUrl     string `json:"source_url" gorm:"column:source_url"`
	Location      string
	MaxPrice      decimal.Decimal
	MinPrice      decimal.Decimal
	Rating        float64 `json:"rating" gorm:"column:rating"`
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

// CreateWithRelations 创建商品和关联数据，并返回细化错误信息
func (productModel *Product) CreateWithRelations() error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("创建商品事务错误: %w", tx.Error)
	}

	baseProduct := *productModel
	baseProduct.Video = nil
	baseProduct.Images = nil
	baseProduct.Skus = nil
	baseProduct.Attrs = nil

	if err := tx.Create(&baseProduct).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建商品错误: %w", err)
	}
	productModel.ID = baseProduct.ID

	if productModel.Video != nil {
		video := *productModel.Video
		video.ProductID = baseProduct.ID
		if err := tx.Create(&video).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建商品视频错误: %w", err)
		}
		productModel.Video = &video
	}

	if len(productModel.Images) > 0 {
		images := make([]productimages.ProductImage, len(productModel.Images))
		for i := range productModel.Images {
			images[i] = productModel.Images[i]
			images[i].ProductID = baseProduct.ID
		}
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建商品图片错误: %w", err)
		}
		productModel.Images = images
	}

	if len(productModel.Skus) > 0 {
		createdSkus := make([]productsku.ProductSku, 0, len(productModel.Skus))
		for i := range productModel.Skus {
			sku := productModel.Skus[i]
			sku.ProductId = baseProduct.ID

			var skuPackage *productskupackage.ProductSkuPackage
			if sku.SkuPackage != nil {
				packageCopy := *sku.SkuPackage
				skuPackage = &packageCopy
			}
			sku.SkuPackage = nil

			if err := tx.Create(&sku).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("创建商品SKU错误: %w", err)
			}

			if skuPackage != nil {
				skuPackage.SkuID = sku.ID
				if err := tx.Create(skuPackage).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("创建商品SKU包裹错误: %w", err)
				}
				sku.SkuPackage = skuPackage
			}

			createdSkus = append(createdSkus, sku)
		}
		productModel.Skus = createdSkus
	}

	if len(productModel.Attrs) > 0 {
		attrs := make([]productattr.ProductAttr, len(productModel.Attrs))
		for i := range productModel.Attrs {
			attrs[i] = productModel.Attrs[i]
			attrs[i].ProductID = baseProduct.ID
		}
		if err := tx.Create(&attrs).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建商品属性错误: %w", err)
		}
		productModel.Attrs = attrs
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交商品事务错误: %w", err)
	}

	return nil
}

// FindByProductID1688 根据 1688 商品ID查找
// 两层处理来避免 cached plan must not change result type：
// 判重查询不再走 SELECT *，改为只查固定字段 id
// 该查询会话禁用预编译语句：Session(&gorm.Session{PrepareStmt: false})
// 这样即使你后续改表结构，也不会因为 products 的列变化触发这个 cached plan 结果类型冲突。
func FindByProductID1688(productID1688 int64) (Product, error) {
	var productModel Product
	err := database.DB.Session(&gorm.Session{PrepareStmt: false}).
		Model(&Product{}).
		Select("id").
		Where("product_id_1688 = ?", productID1688).
		First(&productModel).Error
	if err != nil {
		return Product{}, err
	}
	return productModel, nil
}

// ExistsByProductID1688 检查 1688 商品ID是否已存在
func ExistsByProductID1688(productID1688 int64) (bool, Product, error) {
	productModel, err := FindByProductID1688(productID1688)
	if err == nil {
		return true, productModel, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, Product{}, nil
	}

	return false, Product{}, fmt.Errorf("查询商品错误: %w", err)
}

// GetByIDWithRelations 查询商品并加载关联数据
func GetByIDWithRelations(id uint64) (Product, error) {
	var productModel Product
	err := database.DB.
		Preload("Video").
		Preload("Images").
		Preload("Attrs").
		Preload("Skus").
		Preload("Skus.SkuPackage").
		First(&productModel, id).Error
	if err != nil {
		return Product{}, fmt.Errorf("查询商品详情错误: %w", err)
	}
	return productModel, nil
}

func ListHomeSections(newLimit int, recommendedLimit int) ([]Product, []Product, error) {
	var newProducts []Product
	err := database.DB.
		Where("on_sale = ?", true).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Order("created_at DESC").
		Limit(newLimit).
		Find(&newProducts).Error
	if err != nil {
		return nil, nil, fmt.Errorf("查询新品失败: %w", err)
	}

	excludeIDs := make([]uint64, 0, len(newProducts))
	for i := range newProducts {
		excludeIDs = append(excludeIDs, newProducts[i].ID)
	}

	var recommendedProducts []Product
	query := database.DB.
		Where("on_sale = ?", true).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Order("rating DESC").
		Order("sale_count DESC").
		Order("created_at DESC").
		Limit(recommendedLimit)

	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	err = query.Find(&recommendedProducts).Error
	if err != nil {
		return nil, nil, fmt.Errorf("查询推荐商品失败: %w", err)
	}

	if len(recommendedProducts) == 0 && recommendedLimit > 0 {
		err = database.DB.
			Where("on_sale = ?", true).
			Preload("Images", func(db *gorm.DB) *gorm.DB {
				return db.Order("id ASC")
			}).
			Order("rating DESC").
			Order("sale_count DESC").
			Order("created_at DESC").
			Limit(recommendedLimit).
			Find(&recommendedProducts).Error
		if err != nil {
			return nil, nil, fmt.Errorf("查询推荐商品兜底失败: %w", err)
		}
	}

	return newProducts, recommendedProducts, nil
}
