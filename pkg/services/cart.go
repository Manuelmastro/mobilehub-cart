package services

import (
	"context"
	"errors"

	"mobilehub-cart/pkg/db"
	"mobilehub-cart/pkg/models"
	"mobilehub-cart/pkg/pb"

	productpb "github.com/Manuelmastro/mobilehub-product/pkg/pb"
	"google.golang.org/grpc"
)

type CartServiceServer struct {
	pb.UnimplementedCartServiceServer
	H db.Handler
}

func (s *CartServiceServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	var cartItems []models.CartItem
	if err := s.H.DB.Where("user_id = ?", req.UserId).Find(&cartItems).Error; err != nil {
		return nil, errors.New("failed to fetch cart")
	}

	var response []*pb.CartItem
	for _, item := range cartItems {
		response = append(response, &pb.CartItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Price:       float32(item.Price),
			Quantity:    item.Quantity,
			TotalPrice:  float32(item.Price * float64(item.Quantity)),
		})
	}

	return &pb.GetCartResponse{Items: response}, nil
}

func (s *CartServiceServer) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	// Initialize a ProductService client
	productServiceConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure()) // Adjust address and port as needed
	if err != nil {
		return nil, errors.New("failed to connect to product service")
	}
	defer productServiceConn.Close()

	productClient := productpb.NewProductServiceClient(productServiceConn)

	// Call GetProduct from ProductService
	productResp, err := productClient.GetProduct(ctx, &productpb.GetProductRequest{Id: req.ProductId})
	if err != nil || productResp.Product == nil {
		return nil, errors.New("product not found in product service")
	}

	product := productResp.Product

	// Add or update cart item
	var cartItem models.CartItem
	if err := s.H.DB.Where("user_id = ? AND product_id = ?", req.UserId, req.ProductId).First(&cartItem).Error; err == nil {
		cartItem.Quantity += req.Quantity
	} else {
		cartItem = models.CartItem{
			UserID:      req.UserId,
			ProductID:   req.ProductId,
			ProductName: product.ProductName,
			Price:       float64(product.Price),
			Quantity:    req.Quantity,
		}
	}

	if err := s.H.DB.Save(&cartItem).Error; err != nil {
		return nil, errors.New("failed to add to cart")
	}

	return &pb.AddToCartResponse{Message: "Product added to cart successfully"}, nil
}

func (s *CartServiceServer) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.RemoveFromCartResponse, error) {
	var cartItem models.CartItem
	if err := s.H.DB.Where("user_id = ? AND product_id = ?", req.UserId, req.ProductId).First(&cartItem).Error; err != nil {
		return nil, errors.New("item not found in cart")
	}

	if err := s.H.DB.Delete(&cartItem).Error; err != nil {
		return nil, errors.New("failed to remove item from cart")
	}

	return &pb.RemoveFromCartResponse{Message: "Item removed from cart successfully"}, nil
}
