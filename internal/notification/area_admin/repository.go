package area_admin

import (
	"cst/internal/notification/entity"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	FindOne(condition interface{}, args ...interface{}) (entity.AreaAdmin, error)
	FindAll(condition interface{}, args ...interface{}) ([]entity.AreaAdmin, error)
	Count(condition interface{}, args ...interface{}) (int, error)
	Create(model *entity.AreaAdmin) (*entity.AreaAdmin, error)
	Update(id int, model *entity.AreaAdmin) (*entity.AreaAdmin, error)
	FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.AreaAdmin, error)
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

func (s *repository) FindOne(condition interface{}, args ...interface{}) (entity.AreaAdmin, error) {
	var model entity.AreaAdmin
	err := s.db.Where(condition, args...).First(&model).Error
	if err != nil {
		return entity.AreaAdmin{}, err
	}
	return model, err
}

func (s *repository) FindAll(condition interface{}, args ...interface{}) ([]entity.AreaAdmin, error) {
	var models []entity.AreaAdmin
	err := s.db.Where(condition, args...).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) FindByPagination(condition interface{}, offset, limit int, args ...interface{}) ([]entity.AreaAdmin, error) {
	var models []entity.AreaAdmin
	query := s.db.Where(condition, args...)
	err := query.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (s *repository) Count(condition interface{}, args ...interface{}) (int, error) {
	var count int
	query := s.db.Model(entity.AreaAdmin{}).Where(condition, args...)
	err := query.Count(&count).Error
	return count, err
}

func (s *repository) Create(model *entity.AreaAdmin) (*entity.AreaAdmin, error) {
	err := s.db.Create(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify create error")
	}
	return model, err
}

func (s *repository) Update(id int, model *entity.AreaAdmin) (*entity.AreaAdmin, error) {
	err := s.db.Model(&entity.AreaAdmin{ID: id}).UpdateColumns(&model).Error
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info("notify update error")
	}
	return model, err
}
