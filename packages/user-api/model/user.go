package model

type User struct {
	Base

	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Roles     []string `gorm:"-"`
}
