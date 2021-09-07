package persistence

import (
	"alarm_center/internal/domain/repo"
	"github.com/jinzhu/gorm"
)

// user repo
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository new user repo
func NewUserRepository(db *gorm.DB) repo.UserRepository {
	return &userRepo{db:db}
}

// FindAll find user info
func (u *userRepo) FindAll() ([]*repo.User, error) {
	res := make([]*repo.User, 0)
	//data := make(map[string]interface{})
	err 	:= u.db.Find(&res).Error
	if err 	!= nil {
		return nil, err
	}

	return res, nil
}

func (u *userRepo) Delete(id string) (error) {
	//return userRepo.Delete(id)
	return nil
}
