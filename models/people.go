package models

import (
	"gorm.io/gorm"
)

type People struct {
	gorm.Model
	ID      int64  `json:"id" gorm:"uniqueIndex:idx_people_id,primaryKey"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	CreditCards []CreditCard `gorm:"foreignKey:PeopleId"`
}

type CreditCard struct {
	gorm.Model
	Number   string
	PeopleId int64
}