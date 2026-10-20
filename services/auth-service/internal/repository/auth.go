package repository

import (
	"log"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/model"
	"gorm.io/gorm"
)

type AuthRepo struct{
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *AuthRepo{
	return &AuthRepo{
		db,
	}
}


func (repo *AuthRepo) FindById(id uint) (*model.User, error) {
	var user *model.User;
	err := repo.db.Where("id=?", id).First(&user).Error;
	if err != nil{
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepo) FindByEmail(email string) (*model.User, error) {
	var user *model.User;
	err := repo.db.Where("email=?", email).First(&user).Error;
	if err != nil{
		log.Print(err.Error())
		return nil, err
	}
	return user, nil
}

func (repo *AuthRepo) CreateUser(user *model.User) error {
	if err := repo.db.Create(user).Error; err != nil{
		return err
	}
	return nil
}

func (repo *AuthRepo) UpdateToken(userInput model.User) error {
    return repo.db.
        Model(&model.User{}).
        Where("id = ?", userInput.ID).
        Update("token_version", gorm.Expr("token_version + ?", 1)).
        Error
}