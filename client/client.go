package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	pb "grpc/file"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	CA_CERT     = "cert/ca-cert.pem"
	CA_KEY      = "cert/ca-key.pem"
	CLIENT_CERT = "cert/client-cert.pem"
	CLIENT_KEY  = "cert/client-key.pem"
)

func main() {
	caPem, err := os.ReadFile(CA_CERT)

	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	clientCert, err := tls.LoadX509KeyPair(CLIENT_CERT, CLIENT_KEY)

	if err != nil {
		log.Fatal(err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	tlsCredential := credentials.NewTLS(conf)

	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithTransportCredentials(tlsCredential))

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewDownloadFileServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	filePath := "C:\\Users\\admin\\Downloads\\Win64OpenSSL-3_1_3.exe"

	response, err := client.Download(ctx, &pb.DownloadFileRequest{FilePath: filePath}, grpc.MaxCallRecvMsgSize(math.MaxInt64))

	if err != nil {
		log.Fatal(err)
	}

	fileName, err := saveFile(filePath, response.GetDataBytes())

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File is saved at: .\\%s", fileName)
}

func saveFile(filePath string, dataBytes []byte) (string, error) {
	fileName := filepath.Base(filePath)

	err := os.WriteFile(fileName, dataBytes, os.FileMode(0644))

	if err != nil {
		return "", err
	}

	return fileName, nil
}
