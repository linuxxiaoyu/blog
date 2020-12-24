package user

// User is a table in db
type User struct {
	ID       uint `gorm:"primary_key"`
	Name     string
	Password string
}
