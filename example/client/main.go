package main

import (
	"bytes"
	"flag"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/example/comm"
	"github.com/lucas-clemente/quic-go/h2quic"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

const (
	fullChainCertFile = "fullchain.pem"
	privkeyCertFile   = "privkey.pem"
)

var (
	verbose  bool
	quiet    bool
	certPath string
	repeat   uint
	interval uint
	urls     []string

	logger = utils.DefaultLogger
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose QUIC Connection message")
	flag.BoolVar(&verbose, "verbose", false, "verbose QUIC Connection message")
	flag.BoolVar(&quiet, "q", false, "don't print the data")
	flag.BoolVar(&quiet, "quiet", false, "don't print the data")
	flag.StringVar(&certPath, "cert", ".", "certificate directory")
	flag.UintVar(&repeat, "repeat", 1, "repeat time of the request")
	flag.UintVar(&interval, "interval", 1, "interval of repeat request (seconds)")
	flag.Parse()

	urls = flag.Args()
	logger.SetLogTimeFormat("")
	if verbose {
		logger.SetLogLevel(utils.LogLevelDebug)
	} else {
		logger.SetLogLevel(utils.LogLevelInfo)
	}
}

func main() {
	certFile := filepath.Join(certPath, fullChainCertFile)
	keyFile := filepath.Join(certPath, privkeyCertFile)
	transport := &h2quic.RoundTripper{
		QuicConfig:      &quic.Config{Versions: []protocol.VersionNumber{protocol.Version43}},
		TLSClientConfig: comm.GetTLSConfig(certFile, keyFile),
	}
	client := &http.Client{Transport: transport}
	defer transport.Close()

	for i := uint(0); i < repeat; i++ {
		for _, addr := range urls {
			logger.Infof("\n\n\nGET %s", addr)

			rsp, err := client.Get(addr)
			if err != nil {
				panic(err)
			}
			logger.Infof("Got response for %s: %#v", addr, rsp)

			body := &bytes.Buffer{}
			_, err = io.Copy(body, rsp.Body)
			if err != nil {
				panic(err)
			}
			if quiet {
				logger.Infof("Request Body: %d bytes", body.Len())
			} else {
				logger.Infof("Request Body:")
				logger.Infof("%s", body.Bytes())
			}

			time.Sleep(time.Second * time.Duration(interval))
		}
	}
}
