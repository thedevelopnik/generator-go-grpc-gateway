package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"

	pb "jaxf-github.fanatics.corp/apparel/<%= appname %>/protocol"
)

var (
	pem      []byte
	ok       bool
	key      []byte
	dcreds   credentials.TransportCredentials
	keyPair  *tls.Certificate
	certPool *x509.CertPool
)

type TestSuite struct {
	suite.Suite
}

func TestRESTEndpoints(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	pem, err := ioutil.ReadFile("../certs/test/test.cert.pem")
	if err != nil {
		err = errors.Wrap(err, "failed to create pem")
		grpclog.Fatal(err)
	}

	key, err := ioutil.ReadFile("../certs/test/test.key.pem")
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

	dcreds = credentials.NewTLS(&tls.Config{
		ServerName: "localhost:10000",
		RootCAs:    certPool,
	})
}

func (suite *TestSuite) TestRESTVersion() {
	assert := assert.New(suite.T())

	os.Setenv("VERSION", "1.2.314")
	go Serve(10001, "localhost:", pem, key, keyPair, certPool)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	rs, err := client.Get("https://localhost:10001/api/v1/version")
	assert.Nil(err)
	defer rs.Body.Close()
	bodyBytes, err := ioutil.ReadAll(rs.Body)
	assert.Nil(err)
	assert.Equal("{\"version\":\"1.2.314\"}", string(bodyBytes))
}

func (suite *TestSuite) TestRESTHello() {
	assert := assert.New(suite.T())

	go Serve(10002, "localhost:", pem, key, keyPair, certPool)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	rs, err := client.Get("https://localhost:10002/api/v1/sayhello?name=Bob")
	assert.NotNil(rs)
	assert.Nil(err)
	defer rs.Body.Close()
	bodyBytes, err := ioutil.ReadAll(rs.Body)
	assert.NotNil(bodyBytes)
	assert.Nil(err)
	assert.Equal("{\"message\":\"Hello Bob\"}", string(bodyBytes))
}

func (suite *TestSuite) TestGRPCVersion() {
	os.Setenv("VERSION", "1.2.314")
	go Serve(10003, "localhost:", pem, key, keyPair, certPool)

	conn, err := grpc.Dial("localhost:10003", grpc.WithTransportCredentials(dcreds))
	if err != nil {
		err = errors.Wrap(err, "error with grpc connection initiation")
		grpclog.Fatal(err)
	}
	defer conn.Close()

	em := &pb.Empty{}
	client := pb.NewGreeterClient(conn)
	res, err := client.Version(context.Background(), em)
	if err != nil {
		err = errors.Wrap(err, "error with Version")
		grpclog.Fatal(err)
	}
	assert.Equal(suite.T(), "1.2.314", res.Version)
}

func (suite *TestSuite) TestGRPCHello() {
	go Serve(10004, "localhost:", pem, key, keyPair, certPool)

	conn, err := grpc.Dial("localhost:10004", grpc.WithTransportCredentials(dcreds))
	if err != nil {
		err = errors.Wrap(err, "error with grpc connection initiation")
		grpclog.Fatal(err)
	}
	defer conn.Close()

	hr := &pb.HelloRequest{
		Name: "Bob",
	}
	client := pb.NewGreeterClient(conn)
	res, err := client.SayHello(context.Background(), hr)
	if err != nil {
		err = errors.Wrap(err, "error with SayHello")
		grpclog.Fatal(err)
	}
	assert.Equal(suite.T(), "Hello Bob", res.Message)
}
