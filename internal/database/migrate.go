package database

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Card{},
		&CollectedCard{},
		&Identicon{},
		&Community{},
		&CommunityUser{},
	)
}
