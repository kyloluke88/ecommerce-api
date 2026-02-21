package company

import (
	"api/app/models"
	"api/pkg/database"

	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Company struct {
	models.BaseModel
	Name          string
	Url           string
	CompanyId1688 string `gorm:"column:company_id_1688"`

	models.CommonTimestampsField
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (companyModel *Company) Create() {
	database.DB.Create(&companyModel)
}

// FindOrCreateByCompanyID1688 根据 1688 公司ID查找，不存在则创建
func FindOrCreateByCompanyID1688(companyID1688, name, url string) (Company, error) {
	var companyModel Company

	err := database.DB.Where("company_id_1688 = ?", companyID1688).First(&companyModel).Error
	if err == nil {
		return companyModel, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return Company{}, fmt.Errorf("查询公司错误: %w", err)
	}

	companyModel = Company{
		Name:          name,
		Url:           url,
		CompanyId1688: companyID1688,
	}
	companyErr := database.DB.Create(&companyModel).Error
	if companyErr != nil {
		return Company{}, fmt.Errorf("创建公司错误: %w", companyErr)
	}
	return companyModel, nil
}
