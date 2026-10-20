package service

import (
	"errors"
	"os"
	"time"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/model"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.AuthRepo
}

// Hàm tạo token
func TokenGenerator(user *model.User) (string, string, error){
	accessKey := []byte(os.Getenv("ACCESS_SECRET"))
	refreshKey := []byte(os.Getenv("REFRESH_SECRET"))

	accessClaims := model.TokenClaims{
		UserId: user.ID,
		Role: user.Role,
		TokenVersion: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15*time.Minute)),
		},
	}
	refreshClaims := model.TokenClaims{
		UserId: user.ID,
		Role: user.Role,
		TokenVersion: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*7*time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	)

	access, err := accessToken.SignedString(accessKey)
	refresh, err_rf := refreshToken.SignedString(refreshKey)

	if err != nil {
		return "", "", err
	} 
	if err_rf != nil{
		return "", "", err_rf
	}

	return access, refresh, nil
}

func NewService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{
		repo,
	}
}

func (service *AuthService) Login(email string, password string )(*model.LoginResponse, error){
	user, err := service.repo.FindByEmail(email);
	if err != nil{

		return nil, errors.New("USER_NOT_FOUND")
	}
	
	if cmp_err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); cmp_err != nil{
		return nil, errors.New("WRONG_PASSWORD")
	}

	if upd_err := service.repo.UpdateToken(*user); upd_err != nil{
		return nil, errors.New("TOKEN_VERSION_UPDATE_FAILED")
	}

	user.TokenVersion ++;
	access, refresh, err_tkgen := TokenGenerator(user)

	if err_tkgen != nil{
		return nil, err_tkgen
	}	

	return &model.LoginResponse{
		AccessToken: access,
		RefreshToken: refresh,
	}, nil
}


func (service *AuthService) Register(email string, password string)(error){
	user, _ := service.repo.FindByEmail(email)
	if user != nil{
		return errors.New("EMAIL_ALREADY_USED")
	}

	hashed_password, h_err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if h_err != nil{
		return errors.New("HASH_FAILED")
	}

	return  service.repo.CreateUser(&model.User{
		Email: email,
		Password: string(hashed_password),
	})

}

