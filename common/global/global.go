package global

import (
	"gorm.io/gorm"
)

const (
	// Version is the version of the project.
	Version = "0.1.0"
)

var (
	// gorm DB
	DBAccess *gorm.DB
)
