package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnector interface {
	Connect() apierrors.ApiError
	Migrate(interface{}) apierrors.ApiError
	Add(interface{}) apierrors.ApiError
	FindByID(object interface{}, id string, selectString []string) apierrors.ApiError
	UpdateByObject(data interface{}) apierrors.ApiError
	UpdateByMap(model interface{}, data map[string]interface{}) apierrors.ApiError
	GetAll(data interface{}, selectString []string) apierrors.ApiError
	DeleteByID(model interface{}) apierrors.ApiError
	GetData(data interface{}, selectString []string, id string) apierrors.ApiError
	ExecuteTransaction(data []interface{}) apierrors.ApiError
	FindByWhere(object interface{}, query string, queryData []interface{}) apierrors.ApiError
}

type GormConnector struct {
	url       string
	connector *gorm.DB
}

func NewDatabaseConnector(url string) *GormConnector {
	return &GormConnector{url: url}
}

func (s *GormConnector) Connect() apierrors.ApiError {
	db, err := gorm.Open(sqlite.Open(s.url), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return apierrors.NewGeneralError(err)
		//panic("failed to connect database")
	}
	s.connector = db
	return nil
}

func (s *GormConnector) Migrate(object interface{}) apierrors.ApiError {
	s.connector.AutoMigrate(&object)
	return nil
}

func (s *GormConnector) Add(object interface{}) apierrors.ApiError {
	if dbc := s.connector.Create(object); dbc.Error != nil {
		if strings.Contains(dbc.Error.Error(), "UNIQUE constraint failed:") {
			return apierrors.NewClientError(dbc.Error)
		}
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}
func (s *GormConnector) FindByID(object interface{}, id string, selectString []string) apierrors.ApiError {

	if dbc := s.connector.Select(selectString).First(&object, "id = ?", id); dbc.Error != nil {
		if dbc.Error.Error() == "record not found" {
			return apierrors.NewNotFoundError(dbc.Error)
		}
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}
func (s *GormConnector) UpdateByObject(data interface{}) apierrors.ApiError {
	if dbc := s.connector.Save(data); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}

func (s *GormConnector) UpdateByMap(model interface{}, data map[string]interface{}) apierrors.ApiError {
	if dbc := s.connector.Model(model).Where("id = ?", data["ID"].(string)).Updates(data); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}

func (s *GormConnector) GetAll(data interface{}, selectString []string) apierrors.ApiError {
	if dbc := s.connector.Select(selectString).Find(data); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}

func (s *GormConnector) DeleteByID(data interface{}) apierrors.ApiError {
	if dbc := s.connector.Delete(data); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}

func (s *GormConnector) GetData(data interface{}, selectString []string, id string) apierrors.ApiError {
	if dbc := s.connector.Select(selectString).Where("id = ?", id).Find(data); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}

func (s *GormConnector) ExecuteTransaction(data []interface{}) apierrors.ApiError {
	dbc := s.connector.Transaction(
		func(tx *gorm.DB) error {
			for _, v := range data {
				updateSet := v.([]interface{})
				model := updateSet[0]
				data := updateSet[1].(map[string]interface{})
				query := updateSet[2]
				queryData := updateSet[3].([]interface{})
				dbc := tx.Model(model).Where(query, queryData...).Updates(data)
				if dbc.Error != nil {
					return dbc.Error
				}
				if dbc.RowsAffected != 1 {
					return errors.New("Update failed: Rows affected: " + fmt.Sprintln(dbc.RowsAffected))
				}
			}
			return nil
		})
	if dbc != nil {
		return apierrors.NewGeneralError(dbc)
	}
	return nil
}

func (s *GormConnector) FindByWhere(object interface{}, query string, queryData []interface{}) apierrors.ApiError {

	if dbc := s.connector.Where(query, queryData...).First(&object); dbc.Error != nil {
		return apierrors.NewGeneralError(dbc.Error)
	}
	return nil
}
