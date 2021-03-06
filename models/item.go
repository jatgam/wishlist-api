package models

import (
	"github.com/jinzhu/gorm"

	"github.com/jatgam/wishlist-api/db"
)

type ItemModel struct {
	DefaultModel
	Name       string    `gorm:"column:name;type:varchar(255);not null"`
	URL        string    `gorm:"column:url;type:varchar(255);not null"`
	Reserved   bool      `gorm:"column:reserved;type:tinyint(1);not null;DEFAULT:false"`
	Reserver   UserModel `gorm:"foreignkey:ReserverID"`
	ReserverID *int      `gorm:"column:reserverid;type:integer;DEFAULT:NULL"`
	Rank       int       `gorm:"colunm:rank;type:integer;DEFAULT:NULL"`
}

func (ItemModel) TableName() string {
	return "items1"
}

func ItemDefaultScope(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, url, reserved, rank, createdAt, updatedAt")
}

func ItemReserveScope(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, url, reserved, reserverid, rank, createdAt, updatedAt")
}

func ItemOrderScope(db *gorm.DB) *gorm.DB {
	return db.Order("rank ASC, id DESC")
}

func GetItems(condition interface{}, scopes ...func(*gorm.DB) *gorm.DB) (*[]ItemModel, error) {
	db := db.GetDB()
	if len(scopes) < 1 {
		scopes = append(scopes, ItemDefaultScope)
	}
	var model []ItemModel
	err := db.Scopes(scopes...).Where(condition).Find(&model).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &model, err
}

func FindOneItem(condition interface{}, scopes ...func(*gorm.DB) *gorm.DB) (*ItemModel, error) {
	db := db.GetDB()
	if len(scopes) < 1 {
		scopes = append(scopes, ItemDefaultScope)
	}
	var model ItemModel
	err := db.Scopes(scopes...).Where(condition).First(&model).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &model, err
}

func GetWantedItems() (*[]ItemModel, error) {
	condition := map[string]interface{}{"reserved": false}
	items, err := GetItems(condition, ItemDefaultScope, ItemOrderScope)
	return items, err
}

func GetAllItems() (*[]ItemModel, error) {
	condition := map[string]interface{}{}
	items, err := GetItems(condition, ItemDefaultScope, ItemOrderScope)
	return items, err
}

func GetReservedItems(userID int) (*[]ItemModel, error) {
	condition := map[string]interface{}{"reserverid": userID}
	items, err := GetItems(condition, ItemDefaultScope, ItemOrderScope)
	return items, err
}

func AddItem(name, url string, rank int) error {
	db := db.GetDB()
	newItem := &ItemModel{Name: name, URL: url, Rank: rank}
	err := db.Create(newItem).Error
	return err
}

func DeleteItem(item *ItemModel) error {
	db := db.GetDB()
	err := db.Delete(&item).Error
	return err
}

func UpdateItemWithMap(item *ItemModel, updates map[string]interface{}) error {
	db := db.GetDB()
	err := db.Model(item).Updates(updates).Error
	return err
}
