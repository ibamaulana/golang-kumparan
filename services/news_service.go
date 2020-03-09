package services

import (
	"time"

	"github.com/ibamaulana/golang-kumparan/model"
	"github.com/jinzhu/gorm"
)

type NewsServiceContract interface {
	Create(news *model.News, tx *gorm.DB) error
	Find(ID int) (*model.News, error)
	Get() ([]*model.News, error)
}

type newsContractService struct {
	db *gorm.DB
}

func NewNewsServiceContract(db *gorm.DB) NewsServiceContract {
	return &newsContractService{db}
}

func (srv *newsContractService) Create(news *model.News, tx *gorm.DB) error {
	var err error
	// err = tx.Create(&news).Error
	query := "INSERT INTO news (author, body, created) VALUES (?,?,?)"
	err = tx.Exec(query, news.Author, news.Body, time.Now()).Error
	return err
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

func (srv *newsContractService) Find(ID int) (*model.News, error) {
	news := new(model.News)
	var err error

	err = srv.db.Where("id=?", ID).Find(&news).Error
	if err != nil {
		return nil, err
	}

	return news, nil
}
