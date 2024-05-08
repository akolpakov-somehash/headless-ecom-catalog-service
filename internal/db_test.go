package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCreateProduct(t *testing.T) {
	// given
	productId := uint64(1)
	tests := []struct {
		name    string
		product DbProduct
		setup   func(p *DbProduct) *DbWrapperMock
		wantErr bool
	}{
		{
			name: "Create a new product",
			product: DbProduct{
				Name:        "Test Product",
				Sku:         "test-sku",
				Description: "Test Description",
				Price:       10.0,
				Image:       "test.jpg",
				ID:          productId, //have to set it manually
			},
			setup: func(p *DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Create", p).Return(&gorm.DB{}).Once()
				return dbWrapper
			},
			wantErr: false,
		},
		{
			name:    "Create a new product with an error",
			product: DbProduct{},
			setup: func(p *DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Create", p).Return(&gorm.DB{Error: gorm.ErrInvalidData}).Once()
				return dbWrapper
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbWrapper := tt.setup(&tt.product)
			//when
			id, err := CreateProduct(dbWrapper, &tt.product)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, productId, id)
			}

		})
	}
}

func TestGetProductByID(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name    string
		id      uint64
		product DbProduct
		setup   func(id uint64) *DbWrapperMock
		wantErr bool
	}{
		{
			name:    "Fetch existing product",
			id:      productID,
			product: DbProduct{ID: productID, Name: "Existing Product"},
			setup: func(id uint64) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("First", &DbProduct{}, []interface{}{id}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*DbProduct)
					arg.ID = id
					arg.Name = "Existing Product"
				}).Return(&gorm.DB{}).Once()
				return dbWrapper
			},
			wantErr: false,
		},
		{
			name: "Fetch non-existing product",
			id:   productID,
			setup: func(id uint64) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("First", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return dbWrapper
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			dbWrapper := tt.setup(tt.id)
			product, err := GetProductByID(dbWrapper, tt.id)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, &tt.product, product)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name    string
		product DbProduct
		setup   func(p *DbProduct) *DbWrapperMock
		wantErr bool
	}{
		{
			name:    "Update existing product",
			product: DbProduct{ID: productID, Name: "Updated Product"},
			setup: func(p *DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Save", p).Return(&gorm.DB{}).Once()
				return dbWrapper
			},
			wantErr: false,
		},
		{
			name:    "Update non-existing product",
			product: DbProduct{ID: productID, Name: "Updated Product"},
			setup: func(p *DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Save", p).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return dbWrapper
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			dbWrapper := tt.setup(&tt.product)
			err := UpdateProduct(dbWrapper, &tt.product)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteProductByID(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name      string
		productID uint64
		setup     func(id uint64) *DbWrapperMock
		wantErr   bool
	}{
		{
			name:      "Delete existing product",
			productID: productID,
			setup: func(id uint64) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Delete", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{}).Once()
				return dbWrapper
			},
			wantErr: false,
		},
		{
			name:      "Delete non-existing product",
			productID: productID,
			setup: func(id uint64) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Delete", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return dbWrapper
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			dbWrapper := tt.setup(tt.productID)
			err := DeleteProductByID(dbWrapper, tt.productID)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	// given
	tests := []struct {
		name     string
		products []*DbProduct
		setup    func(ps []*DbProduct) *DbWrapperMock
		wantErr  bool
	}{
		{
			name: "Fetch all products",
			products: []*DbProduct{
				{ID: 1, Name: "Product 1"},
				{ID: 2, Name: "Product 2"},
			},
			setup: func(ps []*DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Find", mock.AnythingOfType("*[]*internal.DbProduct"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]*DbProduct)
					*arg = ps
				}).Return(&gorm.DB{}).Once()
				return dbWrapper
			},
			wantErr: false,
		},
		{
			name: "Fetch all products with an error",
			setup: func(ps []*DbProduct) *DbWrapperMock {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Find", mock.AnythingOfType("*[]*internal.DbProduct"), mock.Anything).Return(&gorm.DB{Error: gorm.ErrInvalidData}).Once()
				return dbWrapper
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			dbWrapper := tt.setup(tt.products)
			products, err := GetAllProducts(dbWrapper)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Len(t, products, len(tt.products))
				assert.Equal(t, tt.products, products)
				assert.NoError(t, err)
			}
		})
	}
}
