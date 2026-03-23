package model

import (
	"gorm.io/gorm"
)

/*
* RegisterTables create all database tables in this function.
* Notice! you should create tables at here!
 */
func RegisterTables(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Bsp{})
	if err != nil {
		return err
	}

	return nil
}
