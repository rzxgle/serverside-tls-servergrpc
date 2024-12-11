package main

import (
	"context"
	"crypto/tls"
	"exemploserversidetls/src/pb/products"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	products.ProductServiceServer
}

func (s *server) FindAll(ctx context.Context, req *products.ListProductRequest) (*products.ListProductResponse, error) {
	productList := make([]*products.Product, 0)
	productList = append(productList, &products.Product{
		Id:    1,
		Title: "Laptop Dell",
	})

	return &products.ListProductResponse{
		Products: productList,
	}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("./src/cert/server-cert.pem", "./src/cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	srv := server{}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	//carregar as credenciais
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalln(err)
	}

	//esse servidor usará tls na criação do servidor grpc
	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	products.RegisterProductServiceServer(s, &srv)

	if err := s.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
