package models

type ItemModel struct {
	DefaultModel
	Name       string    `gorm:"column:name;type:varchar(255);not null"`
	URL        string    `gorm:"column:url;type:varchar(255);not null"`
	Reserved   bool      `gorm:"column:reserved;type:tinyint(1);not null;DEFAULT:false"`
	Reserver   UserModel `gorm:"foreignkey:ReserverID"`
	ReserverID int       `gorm:"column:reserverid;type:integer;DEFAULT:NULL"`
	Rank       int       `gorm:"colunm:rank;type:integer;DEFAULT:NULL"`
}

func (ItemModel) TableName() string {
	return "items1"
}
