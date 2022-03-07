package admin

import (
	"cst/internal/notification/entity"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	FindOne(condition interface{}, args ...interface{}) (entity.Admin, error)
	FindAll(condition interface{}, args ...interface{}) ([]entity.Admin, error)
	Count(condition interface{}, args ...interface{}) (int, error)
	Create(model *entity.Admin) (*entity.Admin, error)
	Update(id int, model *entity.Admin) (*entity.Admin, error)
	FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.Admin, error)
}

type repository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRepository(db *gorm.DB, logger *logrus.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (s *repository) FindOne(condition interface{}, args ...interface{}) (entity.Admin, error) {
	var model entity.Admin
	err := s.db.Where(condition, args...).First(&model).Error
	if err != nil {
		return entity.Admin{}, err
	}
	return model, err
}

func (s *repository) FindAll(condition interface{}, args ...interface{}) ([]entity.Admin, error) {
	var models []entity.Admin
	err := s.db.Where(condition, args...).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.Admin, error) {
	var models []entity.Admin
	query := s.db.Where(condition, args...)
	err := query.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) Count(condition interface{}, args ...interface{}) (int, error) {
	var count int
	query := s.db.Model(entity.Admin{}).Where(condition, args...)
	err := query.Count(&count).Error
	return count, err
}

func (s *repository) Create(model *entity.Admin) (*entity.Admin, error) {
	err := s.db.Create(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify create error")
	}
	return model, err
}

func (s *repository) Update(id int, model *entity.Admin) (*entity.Admin, error) {
	err := s.db.Model(&entity.Admin{ID: id}).UpdateColumns(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify update error")
	}
	return model, err
}
