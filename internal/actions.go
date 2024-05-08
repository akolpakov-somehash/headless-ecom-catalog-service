package internal

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/akolpakov-somehash/headless-ecom-protos/gen/go/catalog"
)

type Server struct {
	ProductService ProductServiceInterface
	pb.UnimplementedProductInfoServer
}

func (s *Server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductId, error) {
	dbProduct := protoToProduct(in)
	id, err := s.ProductService.CreateProduct(dbProduct)
	if err != nil {
		log.Printf("Failed to add product %v : %v. Error: %v", id, in.Name, err)
		return nil, fmt.Errorf("failed to add product: %w", err)
	}
	log.Printf("Product %v : %v - Added.", id, in.Name)
	return &pb.ProductId{Id: id}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Empty, error) {
	if _, exists := s.ProductService.GetProductByID(in.Id); exists != nil {
		log.Printf("Failed to find product %v : %v. Error: %v", in.Id, in.Name, exists)
		return nil, fmt.Errorf("product not found: %w", exists)
	}
	updatedProduct := protoToProduct(in)
	if err := s.ProductService.UpdateProduct(updatedProduct); err != nil {
		log.Printf("Failed to update product %v : %v. Error: %v", in.Id, in.Name, err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	log.Printf("Product %v : %v - Updated.", in.Id, in.Name)
	return new(pb.Empty), nil
}

func (s *Server) DeleteProduct(ctx context.Context, in *pb.ProductId) (*pb.Empty, error) {
	if err := s.ProductService.DeleteProductByID(in.Id); err != nil {
		return nil, err
	}
	return new(pb.Empty), nil
}

func (s *Server) GetProductInfo(ctx context.Context, in *pb.ProductId) (*pb.Product, error) {
	dbProduct, err := s.ProductService.GetProductByID(in.Id)
	if err != nil {
		log.Printf("Failed to find product %v. Error: %v", in.Id, err)
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return productToProto(dbProduct), nil
}

func (s *Server) GetProductList(ctx context.Context, in *pb.Empty) (*pb.ProductList, error) {
	dbProducts, err := s.ProductService.GetAllProducts()
	if err != nil {
		log.Printf("Failed to obtain product list. Error: %v", err)
		return nil, fmt.Errorf("failed to obtain product list: %w", err)
	}
	protoProducts := make(map[uint64]*pb.Product, len(dbProducts))
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, product := range dbProducts {
		wg.Add(1)
		go func(product *DbProduct) {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			protoProducts[product.ID] = productToProto(product)
		}(product)
	}
	wg.Wait()

	return &pb.ProductList{Products: protoProducts}, nil
}

func protoToProduct(product *pb.Product) *DbProduct {
	return &DbProduct{
		ID:          product.Id,
		Name:        product.Name,
		Sku:         product.Sku,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
	}
}

func productToProto(dbProduct *DbProduct) *pb.Product {
	return &pb.Product{
		Id:          dbProduct.ID,
		Name:        dbProduct.Name,
		Sku:         dbProduct.Sku,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
		Image:       dbProduct.Image,
	}
}
