package service

import (
	"errors"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type CoverService interface {
	Service
	Get(int64) (*model.Cover, error)
	Insert(*model.Cover) error
	Update(*model.Cover) error
	Delete(*model.Cover) error
	GetOneBy(string, interface{}) (*model.Cover, error)
	FindAll(string) ([]*model.Cover, error)
}

type coverServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewCoverService(dbMap *gorp.DbMap) CoverService {
	return &coverServiceImpl{dbMap: dbMap}
}

func (s *coverServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *coverServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Cover{}, "covers").SetKeys(true, "id")
	return nil
}

func (s *coverServiceImpl) Get(id int64) (*model.Cover, error) {
	c, err := s.DbMap().Get(model.Cover{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	cc, ok := c.(*model.Cover)
	if !ok {
		return nil, errors.New("Unable to type cast cover")
	}

	if cc.Id == int64(0) {
		return nil, nil
	}

	return cc, nil
}

func (s *coverServiceImpl) Insert(c *model.Cover) error {
	return s.DbMap().Insert(c)
}

func (s *coverServiceImpl) Update(c *model.Cover) error {
	_, err := s.DbMap().Update(c)
	return err
}

func (s *coverServiceImpl) Delete(c *model.Cover) error {
	_, err := s.DbMap().Delete(c)
	return err
}

func (s *coverServiceImpl) GetOneBy(column string, value interface{}) (*model.Cover, error) {
	var cover model.Cover
	err := s.DbMap().SelectOne(&cover, "select * from covers where "+column+" = ? limit 1", value)
	return &cover, err
}

func (s *coverServiceImpl) FindAll(order string) ([]*model.Cover, error) {
	sql := `select * from covers order by ?`

	var covers []*model.Cover
	_, err := s.dbMap.Select(&covers, sql, order)

	return covers, err
}
