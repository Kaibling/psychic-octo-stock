package repositories

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/lucsky/cuid"
)

var DEFAULT_TOKEN_LENGTH = 32

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

func (s *TokenRepository) GenerateAndAddToken(userid string, validUntil int64) (string, apierrors.ApiError) {

	tokenID := cuid.New()
	tokenString := GenerateToken()
	token := &models.Token{ID: tokenID, Active: true, UserID: userid, Token: tokenString, ValidUntil: validUntil}
	if err := s.db.Add(&token); err != nil {
		return "", err
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
func GenerateToken() string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, DEFAULT_TOKEN_LENGTH)
	for i := 0; i < DEFAULT_TOKEN_LENGTH; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			fmt.Println(err)
			return ""
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func (s *TokenRepository) GetUserIDByToken(token string) (string, apierrors.ApiError) {
	var object models.Token
	if err := s.db.FindByWhere(&object, "token = ?", []interface{}{token}); err != nil {
		return "", err
	}
	return object.UserID, nil

}
