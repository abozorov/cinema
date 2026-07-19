package user

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/cache"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/abozorov/cinema/pkg/jwt"
	mailsender "github.com/abozorov/cinema/pkg/mail_sender"
	"github.com/abozorov/cinema/pkg/password"
)

type UserService struct {
	serviceManager services.IServiceManager
	jwt            *jwt.JWTSecret
	memCache       cache.ICache
	mailSender     *mailsender.MailSender
}

func NewUserService(
	serviceManager services.IServiceManager,
	jwt *jwt.JWTSecret,
	memCache cache.ICache,
	mailsender *mailsender.MailSender) *UserService {

	return &UserService{
		serviceManager: serviceManager,
		jwt:            jwt,
		memCache:       memCache,
		mailSender:     mailsender,
	}
}

type user struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password_hash"`
	Age          int    `json:"age"`
}

type sendOtp struct {
	Code       int  `json:"code"`
	User       user `json:"user"`
	AttemptOTP int  `json:"attempt_otp"`
}

func (u *UserService) Verification(ctx context.Context, req models.Verification) (int, error) {
	// check for exist in db
	_, err := u.serviceManager.UserService().GetByEmail(ctx, &userv1.GetUserByEamilRequest{
		Email: req.Email,
	})
	if err == nil {
		return 0, fmt.Errorf("user_service.Verification: %w", err)
	}

	// check mem cash for exist
	var user sendOtp
	// log.Println("verrify req .emain \"", req.Email, "\"")

	err = u.memCache.Get(ctx, req.Email, &user)
	if err != nil {
		return 0, fmt.Errorf("user_service.Verification: %w %w", err, errs.ErrVerifyingFailed)
	}
	user.AttemptOTP++

	if user.AttemptOTP > 2 {
		u.memCache.Delete(ctx, req.Email)
		return 0, fmt.Errorf("user_service.Verification: %w", errs.ErrToManyAttempt)
	}

	if user.Code != req.Code {
		err := u.memCache.Save(ctx, req.Email, user, time.Minute*5)
		if err != nil {
			return 0, fmt.Errorf("user_service.Verification u.memCache.Save: %w", err)
		}
		return 0, fmt.Errorf("user_service.Verification: %w", errs.ErrIncorrectOTPCode)
	}
	u.memCache.Delete(ctx, req.Email)

	userPesponse := user.User
	id, err := u.serviceManager.UserService().Add(ctx, &userv1.CreateUserRequest{
		Name:     userPesponse.Name,
		Email:    userPesponse.Email,
		Phone:    userPesponse.Phone,
		Age:      int32(userPesponse.Age),
		Password: userPesponse.PasswordHash,
	})
	if err != nil {
		return 0, fmt.Errorf("user_service.Register: %w", err)
	}

	return int(id.GetId()), nil
}

func (u *UserService) Register(ctx context.Context, request models.RegisterRequest) error {
	err := request.Validate()
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	// check for exist in db
	_, err = u.serviceManager.UserService().GetByEmail(ctx, &userv1.GetUserByEamilRequest{
		Email: request.Email,
	})
	if err == nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	// check for exist in memcache
	var a struct{}
	err = u.memCache.Get(ctx, request.Email, a)
	if err == nil {
		return fmt.Errorf("user_service.Register: %w %w", err, errs.ErrUserNotBeenVerified)
	}

	// hash password
	request.Password, err = password.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	user := user{
		Name:         request.Name,
		Email:        request.Email,
		Phone:        request.Phone,
		PasswordHash: request.Password,
	}

	otpCode := rand.Int()%899999 + 100000

	// save
	err = u.memCache.Save(ctx, user.Email, sendOtp{
		Code: otpCode,
		User: user,
	}, time.Minute*5)
	if err != nil {
		return fmt.Errorf("user_service.Register u.memCache.Save: %w", err)
	}

	err = u.mailSender.SendMail(user.Email, strconv.Itoa(otpCode))
	if err != nil {
		u.memCache.Delete(ctx, user.Email)
		return fmt.Errorf("user_service.Register: %w", err)
	}

	return nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := u.serviceManager.UserService().GetByID(ctx, &userv1.GetUserRequest{
		Id: int64(id),
	})
	if err != nil {
		return &models.User{}, fmt.Errorf("user_service.GetByID: %w", err)
	}
	return &models.User{
		ID:    int(user.GetId()),
		Name:  user.GetName(),
		Email: user.GetEmail(),
		Phone: user.GetPhone(),
		Age:   int(user.GetAge()),
	}, nil
}

func (u *UserService) Update(ctx context.Context, user models.User) error {
	// validation
	user.Name = strings.TrimSpace(user.Name)
	user.Phone = strings.TrimSpace(user.Phone)

	// updating
	_, err := u.serviceManager.UserService().Update(ctx, &userv1.UpdateUserRequest{
		Id:    int64(user.ID),
		Name:  user.Name,
		Phone: user.Phone,
	})
	if err != nil {
		return fmt.Errorf("user_service.Update: %w", err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, request models.LoginRequest) (*models.Tokens, error) {
	err := request.Validate()
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// get user by email
	user, err := u.serviceManager.UserService().GetByEmail(ctx, &userv1.GetUserByEamilRequest{
		Email: request.Email,
	})
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// compare password
	err = password.Compare(user.GetPassword(), request.Password)
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// generate tokens
	jwtToken, err := u.jwt.GenerateToken(int(user.GetId()), user.GetEmail(), user.GetRole())
	if err != nil {
		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	}

	// refreshToken := refreshtoken.Generate()
	// _, err = u.serviceManager.UserService().GetByID(ctx, &userv1.GetUserRequest{
	// 	Id: user.Id,
	// })
	// if exist, _ := u.refreshTokenR.ExistByUserID(ctx, user.ID); exist {
	// 	err = u.refreshTokenR.DeleteByUserID(ctx, user.ID)
	// 	if err != nil {
	// 		return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	// 	}
	// }
	// err = u.refreshTokenR.Create(ctx, models.RefreshToken{
	// 	UserID:    user.ID,
	// 	TokenHash: refreshtoken.HashRefreshToken(refreshToken),
	// 	ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	// 	CreatedAt: time.Now(),
	// })
	// if err != nil {
	// 	return &models.Tokens{}, fmt.Errorf("user_service.Login: %w", err)
	// }

	// return tokens
	return &models.Tokens{
		// Refresh: refreshToken,
		JWT: jwtToken,
	}, nil
}
