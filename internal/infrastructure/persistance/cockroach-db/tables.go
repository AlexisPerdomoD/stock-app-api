package cockroachdb

import (
	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"gorm.io/gorm"
)

type userTable struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string
	Active   bool
}

func (u *userTable) ToDomain() *domain.User {
	
	return &domain.User{
		ID:       u.ID,
		UserName: u.UserName,
		Password: u.Password,
		Active:   u.Active,
	}

}
