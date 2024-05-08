package internal

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/akolpakov-somehash/headless-ecom-protos/gen/go/catalog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestServer_AddProduct(t *testing.T) {
	// given
	testCases := []struct {
		name        string
		product     *pb.Product
		expectedId  *pb.ProductId
		expectedErr error
		setup       func(p *DbProduct) *ProductServiceMock
	}{
		{
			name: "Create a new product",
			product: &pb.Product{
				Name:        "Test Product",
				Sku:         "test-sku",
				Description: "Test Description",
				Price:       100.0,
				Image:       "test-image",
			},
			expectedId:  &pb.ProductId{Id: 1},
			expectedErr: nil,
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("CreateProduct", p).Return(uint64(1), nil)
				return mockProductService
			},
		},
		{
			name:        "Create a new product with an error",
			product:     &pb.Product{}, // empty product
			expectedId:  nil,
			expectedErr: fmt.Errorf("failed to add product: %w", gorm.ErrInvalidData),
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("CreateProduct", p).Return(uint64(ErrorId), gorm.ErrInvalidData)
				return mockProductService
			},
		},
	}

	for _, tc := range testCases {
		// when
		mockProductService := tc.setup(protoToProduct(tc.product))
		server := &Server{
			ProductService: mockProductService,
		}

		ctx := context.Background()
		productId, err := server.AddProduct(ctx, tc.product)

		// then
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedId, productId)
	}
}

func TestServer_UpdateProductInfo(t *testing.T) {
	// given
	testCases := []struct {
		name           string
		product        *pb.Product
		expecterResult *pb.Empty
		expectedErr    error
		setup          func(p *DbProduct) *ProductServiceMock
	}{
		{
			name: "Update a product",
			product: &pb.Product{
				Id:          1,
				Name:        "Test Product",
				Sku:         "test-sku",
				Description: "Test Description",
				Price:       100.0,
				Image:       "test-image",
			},
			expecterResult: new(pb.Empty),
			expectedErr:    nil,
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetProductByID", p.ID).Return(p, nil)
				mockProductService.On("UpdateProduct", p).Return(nil)
				return mockProductService
			},
		},
		{
			name: "Update a product with an error",
			product: &pb.Product{
				Id: 1,
			},
			expectedErr:    fmt.Errorf("failed to update product: %w", gorm.ErrInvalidData),
			expecterResult: nil,
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetProductByID", p.ID).Return(p, nil)
				mockProductService.On("UpdateProduct", p).Return(gorm.ErrInvalidData)
				return mockProductService
			},
		},
		{
			name: "Update a missing product",
			product: &pb.Product{
				Id: 1,
			},
			expectedErr:    fmt.Errorf("product not found: failed to get a product 1: record not found"),
			expecterResult: nil,
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetProductByID", p.ID).Return(nil, fmt.Errorf("failed to get a product %d: %w", p.ID, gorm.ErrRecordNotFound))
				return mockProductService
			},
		},
	}

	for _, tc := range testCases {
		// when
		mockProductService := tc.setup(protoToProduct(tc.product))
		server := &Server{
			ProductService: mockProductService,
		}

		ctx := context.Background()
		res, err := server.UpdateProduct(ctx, tc.product)

		// then
		if tc.expectedErr != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expecterResult, res)
	}
}

func TestServer_DeleteProduct(t *testing.T) {
	// given
	testCases := []struct {
		name           string
		productId      *pb.ProductId
		expectedResult *pb.Empty
		expectedErr    error
		setup          func(id uint64) *ProductServiceMock
	}{
		{
			name: "Delete a product",
			productId: &pb.ProductId{
				Id: 1,
			},
			expectedResult: new(pb.Empty),
			expectedErr:    nil,
			setup: func(id uint64) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("DeleteProductByID", id).Return(nil)
				return mockProductService
			},
		},
		{
			name: "Delete a product with an error",
			productId: &pb.ProductId{
				Id: 1,
			},
			expectedResult: nil,
			expectedErr:    gorm.ErrRecordNotFound,
			setup: func(id uint64) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("DeleteProductByID", id).Return(gorm.ErrRecordNotFound)
				return mockProductService
			},
		},
	}

	for _, tc := range testCases {
		// when
		mockProductService := tc.setup(tc.productId.Id)
		server := &Server{
			ProductService: mockProductService,
		}

		ctx := context.Background()
		res, err := server.DeleteProduct(ctx, tc.productId)

		// then
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedResult, res)
	}
}

func TestServer_GetProductInfo(t *testing.T) {
	// given
	testCases := []struct {
		name           string
		productId      *pb.ProductId
		expectedResult *pb.Product
		expectedErr    error
		setup          func(p *DbProduct) *ProductServiceMock
	}{
		{
			name: "Get a product",
			productId: &pb.ProductId{
				Id: 1,
			},
			expectedResult: &pb.Product{
				Id:          1,
				Name:        "Test Product",
				Sku:         "test-sku",
				Description: "Test Description",
				Price:       100.0,
				Image:       "test-image",
			},
			expectedErr: nil,
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetProductByID", p.ID).Return(p, nil)
				return mockProductService
			},
		},
		{
			name: "Get a product with an error",
			productId: &pb.ProductId{
				Id: 2,
			},
			expectedResult: nil,
			expectedErr:    fmt.Errorf("product not found: failed to get a product 2: record not found"),
			setup: func(p *DbProduct) *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetProductByID", p.ID).Return(nil, fmt.Errorf("failed to get a product %d: %w", p.ID, gorm.ErrRecordNotFound))
				return mockProductService
			},
		},
	}

	for _, tc := range testCases {
		// when
		var mockProductService *ProductServiceMock
		if tc.expectedResult != nil {
			mockProductService = tc.setup(protoToProduct(tc.expectedResult))
		} else {
			mockProductService = tc.setup(protoToProduct(&pb.Product{Id: tc.productId.Id}))
		}
		server := &Server{
			ProductService: mockProductService,
		}

		ctx := context.Background()
		res, err := server.GetProductInfo(ctx, tc.productId)

		// then
		if tc.expectedErr != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expectedResult, res)
	}
}

func TestServer_GetProductList(t *testing.T) {
	// given
	testCases := []struct {
		name           string
		expectedResult *pb.ProductList
		expectedErr    error
		setup          func() *ProductServiceMock
	}{
		{
			name: "Get a product list",
			expectedResult: &pb.ProductList{
				Products: map[uint64]*pb.Product{
					1: {
						Id:          1,
						Name:        "Test Product",
						Sku:         "test-sku",
						Description: "Test Description",
						Price:       100.0,
						Image:       "test-image",
					},
				},
			},
			expectedErr: nil,
			setup: func() *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetAllProducts").Return([]*DbProduct{{
					ID:          1,
					Name:        "Test Product",
					Sku:         "test-sku",
					Description: "Test Description",
					Price:       100.0,
					Image:       "test-image",
				}}, nil)
				return mockProductService
			},
		},
		{
			name:           "Get a product list with an error",
			expectedResult: nil,
			expectedErr:    fmt.Errorf("failed to obtain product list: %w", gorm.ErrInvalidData),
			setup: func() *ProductServiceMock {
				mockProductService := new(ProductServiceMock)
				mockProductService.On("GetAllProducts").Return(nil, gorm.ErrInvalidData)
				return mockProductService
			},
		},
	}

	for _, tc := range testCases {
		// when
		mockProductService := tc.setup()
		server := &Server{
			ProductService: mockProductService,
		}

		ctx := context.Background()
		res, err := server.GetProductList(ctx, new(pb.Empty))

		// then
		if tc.expectedErr != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error())
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, tc.expectedResult, res)
	}
}
