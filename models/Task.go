package models

import "gorm.io/gorm"

type Task struct {
	Id       uint    `gorm:"primary key;autoIncrement" json:"id"`
	Task     *string `json:"task"`
	Priority *string `json:"priority"`
	Is_done  *bool   `json:"is_done"`
}

func MigrateTasks(db *gorm.DB) error {
	err := db.AutoMigrate(&Task{})
	return err
}
