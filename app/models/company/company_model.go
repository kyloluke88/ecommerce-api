package company

import (
	"api/app/models"
	"api/pkg/database"
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
