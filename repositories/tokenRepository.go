package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/lucsky/cuid"
)

var TokenRepo *TokenRepository

type TokenRepository struct {
	db database.DBConnector
}

func SetTokenRepo(repo *TokenRepository) {
	TokenRepo = repo
}

func NewTokenRepository(dbConn database.DBConnector) *TokenRepository {
	return &TokenRepository{db: dbConn}
}

func (s *TokenRepository) Add(token *models.Token) apierrors.ApiError {
	token.ID = cuid.New()
	token.Active = true
	if err := s.db.Add(&token); err != nil {
		return err
	}
	return nil
}

func (s *TokenRepository) GenerateAndAddToken(userid string, hmaSecret interface{}, validUntil int64) (string, apierrors.ApiError) {

	tokenID := cuid.New()
	token := &models.Token{ID: tokenID, Active: true, UserID: userid}
	if err := s.db.Add(&token); err != nil {
		return "", err
	}
	tokenString, err := GenerateToken(userid, tokenID, hmaSecret, validUntil)
	if err != nil {
		return "", apierrors.NewGeneralError(err)
	}
	return tokenString, nil
}

func (s *TokenRepository) UpdateWithMap(data map[string]interface{}) apierrors.ApiError {
	if err := s.db.UpdateByMap(models.Token{}, data); err != nil {
		return err
	}
	return nil
}

func (s *TokenRepository) GetByID(id string) (*models.Token, apierrors.ApiError) {
	var object models.Token

	if err := s.db.FindByID(&object, id, models.TokenSelect); err != nil {
		return nil, err
	}
	return &object, nil

}

func (s *TokenRepository) GetAll(userID string) ([]*models.Token, apierrors.ApiError) {
	var tokenList []*models.Token

	query := "user_id = ?"
	querData := []interface{}{userID}
	if err := s.db.FindAllByWhere(&tokenList, query, querData); err != nil {
		return nil, err
	}
	return tokenList, nil
}
func (s *TokenRepository) DeleteByObject(data *models.User) apierrors.ApiError {
	if err := s.db.DeleteByObject(data); err != nil {
		return err
	}
	return nil
}
func (s *TokenRepository) DeleteByID(id string) apierrors.ApiError {

	if err := s.db.DeleteByObject(&models.Token{ID: id}); err != nil {
		return err
	}
	return nil
}
func GenerateToken(userid string, id string, hmacSecret interface{}, validUntil int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":         userid,
		"generationdate": time.Now().Unix(),
		"validuntil":     validUntil,
		"id":             id,
	})

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Parse(tokenString string, hmac []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmac, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}
	parsedData := map[string]interface{}{}
	parsedData["userid"] = claims["userid"]
	parsedData["generationdate"] = claims["generationdate"]
	parsedData["validuntil"] = claims["validuntil"]
	parsedData["id"] = claims["id"]
	return parsedData, nil
}
