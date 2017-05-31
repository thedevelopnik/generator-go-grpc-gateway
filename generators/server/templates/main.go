package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"

	"github.com/pkg/errors"
	"google.golang.org/grpc/grpclog"

	"jaxf-github.fanatics.corp/apparel/<%= appname %>/server"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	addr     string
	keyPair  *tls.Certificate
	certPool *x509.CertPool
	certPath string
	keyPath  string
)

func main() {
	flag.Parse()
	flag.StringVar(&addr, "address", "localhost:", "domain of the server")
	flag.StringVar(&certPath, "certPath", "./certs/test/test.cert.pem", "path to ssl cert file")
	flag.StringVar(&keyPath, "keyPath", "./certs/test/test.key.pem", "path to ssl key file")
	pem, err := ioutil.ReadFile(certPath)
	if err != nil {
		err = errors.Wrap(err, "failed to create pem")
		grpclog.Fatal(err)
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		err = errors.Wrap(err, "failed to create key")
		grpclog.Fatal(err)
	}

	pair, err := tls.X509KeyPair(pem, key)
	if err != nil {
		err = errors.Wrap(err, "failed to create key pair")
		grpclog.Fatal(err)
	}

	keyPair = &pair
	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(pem)
	if !ok {
		err = errors.New("failed to append cert from pem")
		grpclog.Fatal(err)
	}
	server.Serve(*port, addr, pem, key, keyPair, certPool)
}
