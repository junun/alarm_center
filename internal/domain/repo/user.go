package repo

type UserRepository interface {
	FindAll() ([]*User, error)
	Delete(id string) error
	//GetByID(id string) (*User, error)
	//Store(u *User) error
	//Update(u *User) error
}

// User user entry
type User struct {
	Id   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"type:varchar(200)"`
	Age  int
}

// TableName userservice table.
func (User) TableName() string {
	return "user"
}

