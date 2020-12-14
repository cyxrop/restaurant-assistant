package common

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/entity"

	. "restaurant-assistant/pkg/entity"
)

// Create token
func createToken(userID uuid.UUID, tokenType string, tokenID uuid.UUID, pairedTokenID uuid.UUID) (*entity.Token, error) {
	var claims jwt.MapClaims
	var expiredTime time.Time

	switch tokenType {
	case entity.TokenTypeAccess:
		expiredTime = time.Now().Add(time.Minute * AccessTokenLifeTime)
		claims = jwt.MapClaims{
			ClaimKeyTokenId:       tokenID.String(),
			ClaimKeyPairedTokenId: pairedTokenID.String(),
			ClaimKeyUserId:        userID.String(),
			ClaimKeyType:          entity.TokenTypeAccess,
			ClaimKeyExpiredAt:     expiredTime.String(),
		}
		fmt.Printf("claims: %+v \n", claims)
	case entity.TokenTypeRefresh:
		expiredTime = time.Now().Add(time.Hour * RefreshTokenLifeTime)
		claims = jwt.MapClaims{
			ClaimKeyTokenId:       tokenID.String(),
			ClaimKeyPairedTokenId: pairedTokenID.String(),
			ClaimKeyUserId:        userID.String(),
			ClaimKeyType:          entity.TokenTypeRefresh,
			ClaimKeyExpiredAt:     expiredTime.String(),
		}
		fmt.Printf("claims: %+v \n", claims)
	default:
		return nil, base.NewInternalError("unknown token type")
	}

	jwtToken, err := createJWTToken(claims)
	if err != nil {
		return nil, err
	}

	token := &entity.Token{
		ID:          tokenID,
		Value:       jwtToken,
		ExpiredTime: expiredTime,
	}

	return token, nil
}

func createJWTToken(claims jwt.MapClaims) (string, error) {
	var err error

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "", err
	}
	return token, nil
}

// TODO: Link refresh and access tokens by claims
func CreateTokenPair(userID uuid.UUID) (*entity.TokenPair, error) {
	accessTokenD := uuid.NewV4()
	refreshTokenID := uuid.NewV4()

	at, err := createToken(
		userID,
		entity.TokenTypeAccess,
		accessTokenD,
		refreshTokenID,
	)
	if err != nil {
		return nil, err
	}

	rt, err := createToken(
		userID,
		entity.TokenTypeRefresh,
		refreshTokenID,
		accessTokenD,
	)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  *at,
		RefreshToken: *rt,
	}, nil
}

// Extract and validate token
type AccessDetails struct {
	TokenID       string
	Type          string
	PairedTokenID string
}

func ExtractBearToken(r *http.Request) (string, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1], nil
	}

	return "", base.NewAuthError("auth token required")

}

func VerifyToken(r *http.Request) (*AccessDetails, error) {
	tokenString, err := ExtractBearToken(r)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, base.NewAuthError("invalid auth token algorithm")
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, base.NewAuthError("invalid auth token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		tokenId, ok := claims[ClaimKeyTokenId].(string)
		if !ok {
			return nil, base.NewAuthError("invalid auth token")
		}

		tokenType, ok := claims[ClaimKeyType].(string)
		if !ok {
			return nil, base.NewAuthError("invalid auth token")
		}

		pairedTokenId, ok := claims[ClaimKeyPairedTokenId].(string)
		if !ok {
			return nil, base.NewAuthError("invalid auth token")
		}

		return &AccessDetails{
			TokenID:       tokenId,
			Type:          tokenType,
			PairedTokenID: pairedTokenId,
		}, nil
	}

	return nil, base.NewAuthError("invalid auth token")
}
