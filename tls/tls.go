package tls

import (

	"google.golang.org/grpc/credentials"
)
  
func GenerateTLSCreds() (credentials.TransportCredentials, error) {
	// Здесь нужно указать полные пути к файлам
	certFile := "tls/ca/server.crt"
	keyFile := "tls/ca/server.key"
 
 
	return credentials.NewServerTLSFromFile(certFile, keyFile)
 }