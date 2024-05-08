package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestTestProductService_CreateProduct(t *testing.T) {
	// given
	productId := uint64(1)
	tests := []struct {
		name    string
		product DbProduct
		setup   func(p *DbProduct) *ProductService
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
			setup: func(p *DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Create", p).Return(&gorm.DB{}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: false,
		},
		{
			name:    "Create a new product with an error",
			product: DbProduct{},
			setup: func(p *DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Create", p).Return(&gorm.DB{Error: gorm.ErrInvalidData}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := tt.setup(&tt.product)
			//when
			id, err := ps.CreateProduct(&tt.product)
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

func TestProductService_GetProductByID(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name    string
		id      uint64
		product DbProduct
		setup   func(id uint64) *ProductService
		wantErr bool
	}{
		{
			name:    "Fetch existing product",
			id:      productID,
			product: DbProduct{ID: productID, Name: "Existing Product"},
			setup: func(id uint64) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("First", &DbProduct{}, []interface{}{id}).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*DbProduct)
					arg.ID = id
					arg.Name = "Existing Product"
				}).Return(&gorm.DB{}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: false,
		},
		{
			name: "Fetch non-existing product",
			id:   productID,
			setup: func(id uint64) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("First", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			ps := tt.setup(tt.id)
			product, err := ps.GetProductByID(tt.id)
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

func TestProductService_UpdateProduct(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name    string
		product DbProduct
		setup   func(p *DbProduct) *ProductService
		wantErr bool
	}{
		{
			name:    "Update existing product",
			product: DbProduct{ID: productID, Name: "Updated Product"},
			setup: func(p *DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Save", p).Return(&gorm.DB{}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: false,
		},
		{
			name:    "Update non-existing product",
			product: DbProduct{ID: productID, Name: "Updated Product"},
			setup: func(p *DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Save", p).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			ps := tt.setup(&tt.product)
			err := ps.UpdateProduct(&tt.product)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_DeleteProductByID(t *testing.T) {
	// given
	productID := uint64(1)
	tests := []struct {
		name      string
		productID uint64
		setup     func(id uint64) *ProductService
		wantErr   bool
	}{
		{
			name:      "Delete existing product",
			productID: productID,
			setup: func(id uint64) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Delete", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: false,
		},
		{
			name:      "Delete non-existing product",
			productID: productID,
			setup: func(id uint64) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Delete", &DbProduct{}, []interface{}{id}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			ps := tt.setup(tt.productID)
			err := ps.DeleteProductByID(tt.productID)
			//then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	// given
	tests := []struct {
		name     string
		products []*DbProduct
		setup    func(ps []*DbProduct) *ProductService
		wantErr  bool
	}{
		{
			name: "Fetch all products",
			products: []*DbProduct{
				{ID: 1, Name: "Product 1"},
				{ID: 2, Name: "Product 2"},
			},
			setup: func(ps []*DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Find", mock.AnythingOfType("*[]*internal.DbProduct"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]*DbProduct)
					*arg = ps
				}).Return(&gorm.DB{}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: false,
		},
		{
			name: "Fetch all products with an error",
			setup: func(ps []*DbProduct) *ProductService {
				dbWrapper := new(DbWrapperMock)
				dbWrapper.On("Find", mock.AnythingOfType("*[]*internal.DbProduct"), mock.Anything).Return(&gorm.DB{Error: gorm.ErrInvalidData}).Once()
				return &ProductService{DB: dbWrapper}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//when
			ps := tt.setup(tt.products)
			products, err := ps.GetAllProducts()
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

func TestDbProduct_TableName(t *testing.T) {
	// given
	dbProduct := DbProduct{}
	//when
	tableName := dbProduct.TableName()
	//then
	assert.Equal(t, "catalog_products", tableName)
}
