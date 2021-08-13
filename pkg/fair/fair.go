package fair

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StreetFair interface {
	Create(model *Model) (*Model, error)
	All(filters map[string]string) ([]Model, error)
	Delete(registry string) error
	Update(model *Model) error
	Get(registry string) (*Model, error)
}

type sf struct {
	db  *gorm.DB
	log *logrus.Logger
}

// Create creates a new street fair
func (s *sf) Create(model *Model) (*Model, error) {
	if model.Registry == "" {
		return nil, ErrInvalidStreetFair
	}
	if r := s.db.Create(model); r.Error != nil {
		s.log.WithField("model", model).
			Errorf("Creating a new street fair: %+v", r.Error)
		return nil, ErrInternal
	}
	return model, nil
}

func (s *sf) All(filters map[string]string) ([]Model, error) {
	for k, v := range filters {
		if v == "" {
			delete(filters, k)
		}
	}
	var models []Model
	if r := s.db.Where(filters).Find(&models); r.Error != nil {
		s.log.WithField("filters", filters).
			Errorf("Getting all street fairs: %+v", r.Error)
		return nil, ErrInternal
	}
	return models, nil
}

func (s *sf) Delete(registry string) error {
	r := s.db.Where("registry = ?", registry).Delete(Model{})
	if r.Error != nil {
		s.log.WithField("registry", registry).
			Errorf("Deleting a street fair: %+v", r.Error)
		return ErrInternal
	} else if r.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *sf) Update(model *Model) error {
	r := s.db.Where("registry = ?", model.Registry).Updates(model)
	if r.Error != nil {
		s.log.WithField("model", model).
			Errorf("Updating a street fair: %+v", r.Error)
		return ErrInternal
	} else if r.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *sf) Get(registry string) (*Model, error) {
	var model Model
	if r := s.db.Where("registry = ?", registry).First(&model); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		s.log.WithField("registry", registry).
			Errorf("Getting a street fair: %+v", r.Error)
		return nil, ErrInternal
	}
	return &model, nil
}

// New returns a instance concret StreetFair implementation
func New(db *gorm.DB, log *logrus.Logger) (StreetFair, error) {
	if err := db.AutoMigrate(&Model{}); err != nil {
		log.Errorf("Migrating StreetFair: %+v", err)
		return nil, err
	}
	return &sf{db: db, log: log}, nil
}
