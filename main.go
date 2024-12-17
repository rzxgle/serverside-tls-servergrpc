package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"exemploserversidetls/src/pb/products"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	products.ProductServiceServer
}

func (s *server) FindAll(ctx context.Context, req *products.ListProductRequest) (*products.ListProductResponse, error) {
	time.Sleep(2 * time.Second)

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
	pemClientCA, err := os.ReadFile("./src/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("erro ao adicionar o certificado da CA")
	}

	serverCert, err := tls.LoadX509KeyPair("./src/cert/server-cert.pem", "./src/cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("function: ", info.FullMethod)
	return handler(ctx, req)
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
		grpc.UnaryInterceptor(interceptor),
	)

	products.RegisterProductServiceServer(s, &srv)

	if err := s.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
