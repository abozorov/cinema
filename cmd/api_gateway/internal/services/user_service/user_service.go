package user

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/abozorov/cinema/cmd/api_gateway/internal/models"
	"github.com/abozorov/cinema/cmd/api_gateway/internal/services"
	userv1 "github.com/abozorov/cinema/grpc_api/generate/userpb/user/v1"
	"github.com/abozorov/cinema/pkg/errs"
	"github.com/abozorov/cinema/pkg/jwt"
	mailsender "github.com/abozorov/cinema/pkg/mail_sender"
	"github.com/abozorov/cinema/pkg/password"
	"github.com/patrickmn/go-cache"
)

type UserService struct {
	serviceManager services.IServiceManager
	jwt            *jwt.JWTSecret
	memCache       *cache.Cache
	mailSender     *mailsender.MailSender
}

func NewUserService(
	serviceManager services.IServiceManager,
	jwt *jwt.JWTSecret,
	memCache *cache.Cache,
	mailsender *mailsender.MailSender) *UserService {

	return &UserService{
		serviceManager: serviceManager,
		jwt:            jwt,
		memCache:       memCache,
		mailSender:     mailsender,
	}
}

type sendOtp struct {
	code       int
	user       *models.User
	attemptOTP *int
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
	user, ok := u.memCache.Get(req.Email)
	if !ok {
		return 0, fmt.Errorf("user_service.Verification: %w", errs.ErrVerifyingFailed)
	}
	defer func() {
		*user.(sendOtp).attemptOTP++
	}()

	if *user.(sendOtp).attemptOTP > 2 {
		u.memCache.Delete(req.Email)
		return 0, fmt.Errorf("user_service.Verification: %w", errs.ErrToManyAttempt)
	}

	if user.(sendOtp).code != req.Code {
		return 0, fmt.Errorf("user_service.Verification: %w", errs.ErrIncorrectOTPCode)
	}

	userPesponse := *user.(sendOtp).user
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
	u.memCache.Delete(req.Email)

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
	_, ok := u.memCache.Get(request.Email)
	if ok {
		return fmt.Errorf("user_service.Register: %w", errs.ErrUserNotBeenVerified)
	}

	// hash password
	request.Password, err = password.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("user_service.Register: %w", err)
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		Phone:        request.Phone,
		PasswordHash: request.Password,
	}

	otpCode := rand.Int()%899999 + 100000
	attempt := 0

	u.memCache.Set(user.Email, sendOtp{
		code:       otpCode,
		user:       &user,
		attemptOTP: &attempt,
	}, cache.DefaultExpiration)

	err = u.mailSender.SendMail(user.Email, strconv.Itoa(otpCode))
	if err != nil {
		u.memCache.Delete(user.Email)
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
