package productattr

import "api/app/models"

type ProductAttr struct {
	models.BaseModel
	ProductID uint64
	Fid       int64
	Name      string
	Value     string
	IsActive  bool

	models.CommonTimestampsField
}
