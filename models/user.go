package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/jatgam/wishlist-api/db"
	"github.com/jatgam/wishlist-api/utils"
)

// UserModel is the db structure for users
type UserModel struct {
	DefaultModel
	Username             string      `gorm:"column:username;type:varchar(255);unique;not null"`
	PasswordHash         string      `gorm:"column:hash;type:varchar(255);not null"`
	PasswordReset        bool        `gorm:"column:passwordreset;type:tinyint(1);not null;DEFAULT:false"`
	PasswordResetToken   *string     `gorm:"column:passwordResetToken;type:varchar(255);DEFAULT:NULL"`
	PasswordResetExpires *time.Time  `gorm:"column:passwordResetExpires;type:DATETIME;DEFAULT:NULL"`
	UserLevel            uint        `gorm:"column:userlevel;type:tinyint unsigned;not null"`
	EMail                string      `gorm:"column:email;type:varchar(255);not null"`
	FirstName            string      `gorm:"column:firstname;type:varchar(255);not null"`
	LastName             string      `gorm:"column:lastname;type:varchar(255);not null"`
	ReservedItems        []ItemModel `gorm:"foreignkey:ReserverID"`
}

func (UserModel) TableName() string {
	return "users1"
}

func UserDefaultScope(db *gorm.DB) *gorm.DB {
	return db.Select("id, username, userlevel, email, firstname, lastname, createdAt, updatedAt")
}

func UserPassResetScope(db *gorm.DB) *gorm.DB {
	return db.Select("id, username, passwordreset, passwordResetToken, passwordResetExpires, userlevel, email, firstname, lastname, createdAt, updatedAt")
}

func UserAuthScope(db *gorm.DB) *gorm.DB {
	return db.Select("*")
}

func (u *UserModel) SetPassword(password string) error {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *UserModel) checkPassword(password string) error {
	return utils.CheckPassword(password, u.PasswordHash)
}

func (u *UserModel) ValidatePassword(password string) bool {
	err := u.checkPassword(password)
	if err != nil {
		return false
	}
	return true
}

// FindOneUser will search for a user that matches the supplied condition.
// Will mask Record Not Found Errors.
func FindOneUser(condition interface{}, scopes ...func(*gorm.DB) *gorm.DB) (*UserModel, error) {
	db := db.GetDB()
	if len(scopes) < 1 {
		scopes = append(scopes, UserDefaultScope)
	}
	var model UserModel
	err := db.Scopes(scopes...).Where(condition).First(&model).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &model, err
}

func CreateUser(newUser *UserModel) error {
	db := db.GetDB()
	err := db.Create(newUser).Error
	return err
}

func UpdateUser(user *UserModel, updates UserModel) error {
	db := db.GetDB()
	err := db.Model(user).Updates(updates).Error
	return err
}

// UpdateUserWithMap will update an existing user. A map is needed instead
// of the UserModel when you want/need to insert NULL sql field values.
func UpdateUserWithMap(user *UserModel, updates map[string]interface{}) error {
	db := db.GetDB()
	err := db.Model(user).Updates(updates).Error
	return err
}
