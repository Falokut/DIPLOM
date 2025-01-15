package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"dish_as_a_service/conf"
	"dish_as_a_service/domain"
	"dish_as_a_service/entity"
	"dish_as_a_service/jwt"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type AuthRepo interface {
	GetUserByTelegramId(ctx context.Context, telegramId int64) (entity.User, error)
}

type Auth struct {
	cfg              conf.Auth
	telegramBotToken string
	repo             AuthRepo
}

func NewAuth(cfg conf.Auth, telegramBotToken string, repo AuthRepo) Auth {
	return Auth{
		cfg:              cfg,
		telegramBotToken: telegramBotToken,
		repo:             repo,
	}
}

func (s Auth) LoginByTelegram(ctx context.Context, req domain.LoginByTelegramRequest) (*domain.LoginResponse, error) {
	params, err := url.ParseQuery(req.InitTelegramData)
	if err != nil {
		return nil, domain.ErrInvalidTelegramCredentials
	}

	err = verifyTelegramInitData(params, s.telegramBotToken, s.cfg.TelegramExpireDuration)
	if err != nil {
		return nil, errors.WithMessage(domain.ErrInvalidTelegramCredentials, err.Error())
	}

	telegramUserId, err := getUserIdFromTelegramQuery(params)
	if err != nil {
		return nil, errors.WithMessage(err, "get user id from telegram query")
	}

	user, err := s.repo.GetUserByTelegramId(ctx, telegramUserId)
	if err != nil {
		return nil, errors.WithMessage(err, "get user by telegram id")
	}

	roleName := domain.AdminType
	if !user.Admin {
		roleName = domain.UserType
	}
	tokenValue := entity.TokenUserInfo{
		UserId:   user.Id,
		RoleName: roleName,
	}
	accessToken, err := jwt.GenerateToken(
		s.cfg.Access.Secret,
		s.cfg.Access.TTL,
		tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate access token")
	}

	refreshToken, err := jwt.GenerateToken(
		s.cfg.Refresh.Secret,
		s.cfg.Refresh.TTL,
		tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate refresh token")
	}

	return &domain.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s Auth) RefreshAccessToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	tokenValue, err := jwt.ParseToken(refreshToken, s.cfg.Refresh.Secret)
	if err != nil {
		return nil, errors.WithMessage(err, "parse token")
	}
	accessToken, err := jwt.GenerateToken(
		s.cfg.Access.Secret,
		s.cfg.Access.TTL,
		*tokenValue,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "generate access token")
	}
	return accessToken, nil
}

func (s Auth) HasAdminPrivileges(ctx context.Context, accessToken string) (*domain.HasAdminPrivilegesResponse, error) {
	tokenValue, err := jwt.ParseToken(accessToken, s.cfg.Access.Secret)
	if err != nil {
		return nil, domain.DomainInvalidTokenError(err)
	}

	return &domain.HasAdminPrivilegesResponse{
		HasAdminPrivileges: tokenValue.RoleName == domain.AdminType,
	}, nil
}

func getUserIdFromTelegramQuery(q url.Values) (int64, error) {
	userStr := q.Get("user")
	if userStr == "" {
		return -1, domain.ErrInvalidTelegramCredentials
	}
	idStr := gjson.Get(userStr, "id").String()
	if idStr == "" {
		return -1, domain.ErrInvalidTelegramCredentials
	}

	telegramUserId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, domain.ErrInvalidTelegramCredentials
	}
	return telegramUserId, nil
}

func verifyTelegramInitData(q url.Values, token string, expIn time.Duration) error {
	var (
		authDate time.Time
		hash     string
		pairs    = make([]string, 0, len(q))
	)

	for k, v := range q {
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		pairs = append(pairs, k+"="+v[0])
	}

	if hash == "" {
		return domain.ErrTelegramSignMissing
	}
	if expIn > 0 {
		if authDate.IsZero() {
			return domain.ErrTelegramAuthDateMissing
		}

		if authDate.Add(expIn).Before(time.Now()) {
			return domain.ErrTelegramCredentialsExpired
		}
	}

	sort.Strings(pairs)

	if sign(strings.Join(pairs, "\n"), token) != hash {
		return domain.ErrInvalidTelegramCredentials
	}
	return nil
}

func sign(payload, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}
