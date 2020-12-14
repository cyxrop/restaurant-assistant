package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/common"
	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/cache"
	"restaurant-assistant/pkg/storage/database/repository"
)

// TODO: make Service Interfaces
type AuthenticationService struct {
	log *zap.Logger

	// TODO: replace to data provider
	ur *repository.UserRepository
	cp *cache.Provider
}

func NewAuthenticationService(
	log *zap.Logger,
	cp *cache.Provider,
	ur *repository.UserRepository,
) (*AuthenticationService, error) {
	if log == nil {
		return nil, fmt.Errorf("uninitialized logger")
	}

	if cp == nil {
		return nil, fmt.Errorf("uninitialized cache provider")
	}

	if ur == nil {
		return nil, fmt.Errorf("uninitialized user repository")
	}

	return &AuthenticationService{
		log: log,
		cp:  cp,
		ur:  ur,
	}, nil
}

func (as *AuthenticationService) Login(ctx context.Context, username string, password string) (*entity.TokenPair, error) {
	if username == "" {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"empty username",
		)
	}
	if password == "" {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"empty password",
		)
	}

	user, err := as.ur.FindByName(username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("sql.ErrNoRows \n")
			return nil, base.NewError(
				base.ErrInvalidParameters,
				"user with specified username was not found",
			)
		}
		return nil, err
	}

	if !common.CheckPasswordHash(password, user.Password) {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"user with specified parameters was not found",
		)
	}

	tokensPair, err := as.CreateTokenPair(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return tokensPair, nil
}

func (as *AuthenticationService) RefreshToken(ctx context.Context) (*entity.TokenPair, error) {
	tokenID, tokenIDFound := ctx.Value(entity.ContextTokenIDKey).(string)
	pairedTokenID, pairedTokenIDFound := ctx.Value(entity.ContextPairedTokenIdKey).(string)
	tokenType, tokenTypeFound := ctx.Value(entity.ContextTokenTypeKey).(string)

	if !tokenIDFound || !tokenTypeFound || !pairedTokenIDFound {
		return nil, base.NewInternalError("token data not found in ctx")
	}

	if tokenType != entity.TokenTypeRefresh {
		return nil, base.NewError(
			base.ErrInvalidAuthTokenType,
			"invalid token type",
		)
	}

	userID, err := as.cp.Client.Get(ctx, tokenID).Result()
	if err != nil {
		return nil, base.NewError(
			base.ErrAuthTokenExpired,
			"token expired",
		)
	}
	fmt.Printf("userID data: %+v \n", userID)

	userUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, base.NewInternalError(fmt.Sprintf("cannot get uuid from userID string: %+v \n", userID))
	}

	tokensPair, err := as.CreateTokenPair(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	err = as.cp.Client.Del(ctx, tokenID, pairedTokenID).Err()
	if err != nil {
		return nil, base.NewInternalError("cannot delete previous token from cache")
	}

	return tokensPair, nil
}

func (as *AuthenticationService) Logout(ctx context.Context) error {
	tokenID, tokenIDFound := ctx.Value(entity.ContextTokenIDKey).(string)
	pairedTokenID, pairedTokenIDFound := ctx.Value(entity.ContextPairedTokenIdKey).(string)

	if !tokenIDFound || !pairedTokenIDFound {
		return base.NewInternalError("token data not found in ctx")
	}

	exists, err := as.cp.Client.Exists(ctx, tokenID).Result()
	if err != nil {
		return err
	}

	if exists != 1 {
		return base.NewInternalError("token data not found in ctx")
	}

	err = as.cp.Client.XDel(ctx, tokenID, pairedTokenID).Err()
	if err != nil {
		return base.NewInternalError("cannot delete token from cache")
	}

	return nil
}

func (as *AuthenticationService) GetUserFromContext(ctx context.Context) (*entity.User, error) {
	tokenID, tokenIDFound := ctx.Value(entity.ContextTokenIDKey).(string)
	tokenType, tokenTypeFound := ctx.Value(entity.ContextTokenTypeKey).(string)

	if !tokenIDFound || !tokenTypeFound || tokenID == "" || tokenType == "" {
		return nil, base.NewInternalError("token data not found in ctx")
	}

	if tokenType != entity.TokenTypeAccess {
		return nil, base.NewError(
			base.ErrInvalidAuthTokenType,
			"invalid token type",
		)
	}

	userId, err := as.cp.Client.Get(ctx, tokenID).Result()
	if err != nil {
		return nil, base.NewError(
			base.ErrAuthTokenExpired,
			"token expired",
		)
	}

	user, err := as.ur.FindByID(userId)
	if err != nil {
		return nil, base.NewError(
			base.ErrAuthentication,
			"user by specified id not found",
		)
	}

	return user, nil
}

func (as *AuthenticationService) CreateTokenPair(ctx context.Context, userID uuid.UUID) (*entity.TokenPair, error) {
	tokensPair, err := common.CreateTokenPair(userID)
	if err != nil {
		return nil, err
	}

	// Save tokens to cache
	err = as.cp.Client.Set(
		ctx,
		tokensPair.AccessToken.ID.String(),
		userID.String(),
		tokensPair.AccessToken.ExpiredTime.Sub(time.Now()),
	).Err()
	if err != nil {
		return nil, err
	}

	err = as.cp.Client.Set(
		ctx,
		tokensPair.RefreshToken.ID.String(),
		userID.String(),
		tokensPair.RefreshToken.ExpiredTime.Sub(time.Now()),
	).Err()
	if err != nil {
		return nil, err
	}

	return tokensPair, nil
}
