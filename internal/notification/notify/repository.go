package notify

import (
	"cst/internal/notification/entity"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	FindOne(condition interface{}, args ...interface{}) (entity.Notify, error)
	FindAll(condition interface{}, args ...interface{}) ([]entity.Notify, error)
	Count(condition interface{}, args ...interface{}) (int, error)
	Create(model *entity.Notify) (*entity.Notify, error)
	Update(id int, model *entity.Notify) (*entity.Notify, error)
	FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.Notify, error)
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

func (s *repository) FindOne(condition interface{}, args ...interface{}) (entity.Notify, error) {
	var model entity.Notify
	err := s.db.Where(condition, args...).First(&model).Error
	if err != nil {
		return entity.Notify{}, err
	}
	return model, err
}

func (s *repository) FindAll(condition interface{}, args ...interface{}) ([]entity.Notify, error) {
	var models []entity.Notify
	err := s.db.Where(condition, args...).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.Notify, error) {
	var models []entity.Notify
	query := s.db.Where(condition, args...)
	err := query.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) Count(condition interface{}, args ...interface{}) (int, error) {
	var count int
	query := s.db.Model(entity.Notify{}).Where(condition, args...)
	err := query.Count(&count).Error
	return count, err
}

func (s *repository) Create(model *entity.Notify) (*entity.Notify, error) {
	err := s.db.Create(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify create error")
	}
	return model, err
}

func (s *repository) Update(id int, model *entity.Notify) (*entity.Notify, error) {
	err := s.db.Model(&entity.Notify{ID: id}).UpdateColumns(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify update error")
	}
	return model, err
}
