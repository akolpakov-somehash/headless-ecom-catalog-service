package internal

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DbWrapperMock struct {
	mock.Mock
}

func (d *DbWrapperMock) Create(value interface{}) *gorm.DB {
	args := d.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (d *DbWrapperMock) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := d.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (d *DbWrapperMock) Save(value interface{}) *gorm.DB {
	args := d.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (d *DbWrapperMock) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := d.Called(value, conds)
	return args.Get(0).(*gorm.DB)
}

func (d *DbWrapperMock) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := d.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

type ProductServiceMock struct {
	mock.Mock
}

func (p *ProductServiceMock) CreateProduct(product *DbProduct) (uint64, error) {
	args := p.Called(product)
	return args.Get(0).(uint64), args.Error(1)
}

func (p *ProductServiceMock) GetProductByID(id uint64) (*DbProduct, error) {
	args := p.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DbProduct), args.Error(1)
}

func (p *ProductServiceMock) UpdateProduct(product *DbProduct) error {
	args := p.Called(product)
	return args.Error(0)
}

func (p *ProductServiceMock) DeleteProductByID(id uint64) error {
	args := p.Called(id)
	return args.Error(0)
}

func (p *ProductServiceMock) GetAllProducts() ([]*DbProduct, error) {
	args := p.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*DbProduct), args.Error(1)
}
