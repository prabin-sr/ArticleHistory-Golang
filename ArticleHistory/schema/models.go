package schema

import "time"

// User Model definition - User Defined Model
type User struct {
	ID           uint   `gorm:"column:id;primary_key"`
	FirstName    string `gorm:"column:firstname;not null"`
	LastName     string `gorm:"column:lastname;not null"`
	UserName     string `gorm:"column:username;not null"`
	Email        string `gorm:"column:email;unique_index"`
	PasswordHash string `gorm:"column:password_hash;not null"`
	IsActive     bool   `gorm:"column:is_active;default:true;not null"`
	IsAdmin      bool   `gorm:"column:is_admin;default:false;not null"`
	IsDeleted    bool   `gorm:"column:is_deleted;default:false;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// Archive Model definition for Saving meta data about downloaded archives
type Archive struct {
	ID        int  `gorm:"column:id;primary_key"`
	Year      int  `gorm:"column:year;not null"`
	Month     int  `gorm:"column:month;not null"`
	Status    bool `gorm:"column:status;default:false;not null"`
	IsCleaned bool `gorm:"column:is_cleaned;default:false;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
