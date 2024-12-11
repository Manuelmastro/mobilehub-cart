package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Manuelmastro/mobilehub-cart/pkg/config"
	"github.com/Manuelmastro/mobilehub-cart/pkg/db"
	pb "github.com/Manuelmastro/mobilehub-cart/pkg/pb"
	services "github.com/Manuelmastro/mobilehub-cart/pkg/services"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Cart Svc on", c.Port)

	s := services.CartServiceServer{
		H: h,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCartServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
