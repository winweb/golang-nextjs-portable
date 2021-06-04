package models

import "gorm.io/gorm"

type People struct {
	gorm.Model
	Id      int64  `json:"id" gorm:"primaryKey,index:people_id_index,"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}