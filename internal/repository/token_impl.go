package repository

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/vadim-rm/bmstu-web-backend/internal/domain"
	"github.com/vadim-rm/bmstu-web-backend/internal/dto"
	"time"
)

const redisKeyPrefix = "jwt."

type claims struct {
	jwt.StandardClaims
	UserId      uint `json:"userId"`
	IsModerator bool `json:"isModerator"`
}

type TokenConfig struct {
	ExpiresIn time.Duration
	Token     string
	Issuer    string
}

type TokenImpl struct {
	config TokenConfig
	client *redis.Client
}

func NewTokenImpl(config TokenConfig, client *redis.Client) *TokenImpl {
	return &TokenImpl{
		config: config,
		client: client,
	}
}

func (r *TokenImpl) Create(_ context.Context, input dto.TokenClaims) (dto.Token, error) {
	expiresAt := time.Now().Add(r.config.ExpiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    r.config.Issuer,
		},
		UserId:      uint(input.UserId),
		IsModerator: input.IsModerator,
	})
	if token == nil {
		return dto.Token{}, domain.ErrUnknown
	}

	tokenString, err := token.SignedString([]byte(r.config.Token))
	if err != nil {
		return dto.Token{}, fmt.Errorf("error signing token: %w", err)
	}

	return dto.Token{
		ExpiresAt: expiresAt,
		Token:     tokenString,
	}, nil
}

func (r *TokenImpl) GetClaims(ctx context.Context, token string) (dto.TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.config.Token), nil
	})
	if err != nil {
		return dto.TokenClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	blacklisted, err := r.getBlacklisted(ctx, token)
	if err != nil {
		return dto.TokenClaims{}, fmt.Errorf("error getting if token blacklisted: %w", err)
	}

	if blacklisted {
		return dto.TokenClaims{}, fmt.Errorf("found token in blacklist: %w", domain.ErrUnauthenticated)
	}

	tokenClaims := parsedToken.Claims.(*claims)
	return dto.TokenClaims{
		UserId:      domain.UserId(tokenClaims.UserId),
		IsModerator: tokenClaims.IsModerator,
	}, nil
}

func (r *TokenImpl) Blacklist(ctx context.Context, token string) error {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.config.Token), nil
	})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	tokenClaims := parsedToken.Claims.(*claims)
	ttl := time.Now().Sub(time.Unix(tokenClaims.ExpiresAt, 0))
	return r.client.Set(ctx, getJWTKey(token), true, ttl).Err()
}

func (r *TokenImpl) getBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.client.Exists(ctx, getJWTKey(token)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func getJWTKey(token string) string {
	return redisKeyPrefix + token
}
