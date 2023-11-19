package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	pb "grpc/file"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	CA_CERT     = "cert/ca-cert.pem"
	CA_KEY      = "cert/ca-key.pem"
	SERVER_CERT = "cert/server-cert.pem"
	SERVER_KEY  = "cert/server-key.pem"
)

type downloadFileService struct {
	pb.UnimplementedDownloadFileServiceServer
}

func (d *downloadFileService) Download(ctx context.Context, request *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	dataBytes, err := os.ReadFile(request.GetFilePath())

	if err != nil {
		return nil, err
	}

	return &pb.DownloadFileResponse{DataBytes: dataBytes}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:9000")

	if err != nil {
		log.Fatalf("listening port error: %v", err)
	}

	caPem, err := os.ReadFile(CA_CERT)

	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	serverCert, err := tls.LoadX509KeyPair(SERVER_CERT, SERVER_KEY)

	if err != nil {
		log.Fatal(err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	tlsCredentials := credentials.NewTLS(conf)

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))

	pb.RegisterDownloadFileServiceServer(grpcServer, &downloadFileService{})

	log.Printf("listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc serve error: %v", err)
	}
}
