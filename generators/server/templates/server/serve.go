package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "jaxf-github.fanatics.corp/apparel/<%= appname %>/protocol"
)

func grpcHandler(grpcServer *grpc.Server, restHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			restHandler.ServeHTTP(w, r)
		}
	})
}

// Serve is the main function that provides the http and grpc servers
func Serve(port int, addr string, pem []byte, key []byte, keyPair *tls.Certificate, certPool *x509.CertPool) {
	host := fmt.Sprintf("%v%d", addr, port)

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewClientTLSFromCert(certPool, host))}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(grpcServer, newServer())

	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: host,
		RootCAs:    certPool,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux()
	err := pb.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, host, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
	mux.Handle("/", gwmux)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandler(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keyPair},
			NextProtos:   []string{"h2"},
		},
	}

	grpclog.Printf("grpc on port %d\n", port)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		err = errors.Wrap(err, "Fatal error on listen and serve")
		grpclog.Fatal(err)
	}

	return
}
