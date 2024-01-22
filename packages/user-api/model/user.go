package model

type User struct {
	Base

	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Roles     []string `gorm:"-"`
}
