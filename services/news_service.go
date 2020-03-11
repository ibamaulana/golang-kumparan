package services

import (
	"time"

	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/jinzhu/gorm"
)

type NewsServiceContract interface {
	Create(news *model.News, tx *gorm.DB) (*model.News, error)
	Get() ([]*model.News, error)
}

type newsContractService struct {
	db *gorm.DB
}

func NewNewsServiceContract(db *gorm.DB) NewsServiceContract {
	return &newsContractService{db}
}

func (srv *newsContractService) Create(news *model.News, tx *gorm.DB) (*model.News, error) {
	row := new(model.News)
	news.Created = time.Now()
	d := tx.Create(&news).Scan(&row)
	if d.Error != nil {
		return nil, d.Error
	}

	return row, nil
}

func (srv *newsContractService) Get() ([]*model.News, error) {
	var news []*model.News
	var err error

	err = srv.db.Find(&news).Error
	if err != nil {
		return nil, err
	}

	return news, nil
}
