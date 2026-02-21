// Package user 存放用户 Model 相关逻辑
package user

import (
	"api/app/models"
	"api/pkg/database"
	"api/pkg/hash"
	"time"
)

// User 用户模型
type User struct {
	models.BaseModel
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	UserProfile  UserProfile `gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	models.CommonTimestampsField
}

type UserProfile struct {
	models.BaseModel
	UserID        uint64
	FirstName     string
	LastName      string
	FirstNameKana *string
	LastNameKana  *string
	Phone         *string
	DateOfBirth   *time.Time
	Gender        *int
	models.CommonTimestampsField
}

func (User) TableName() string {
	return "users"
}

// Create 创建用户，通过 User.ID 来判断是否创建成功
func (userModel *User) Create() {
	database.DB.Create(&userModel)
}

// ComparePassword 密码是否正确
func (userModel *User) ComparePassword(_password string) bool {
	return hash.BcryptCheck(_password, userModel.PasswordHash)
}

func (userModel *User) Save() (rowsAffected int64) {
	result := database.DB.Save(&userModel)
	return result.RowsAffected
}
