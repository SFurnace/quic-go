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
	verbose   bool
	quiet     bool
	certPath  string
	repeat    uint
	interval  uint
	timeout   uint
	keepAlive bool
	urls      []string

	logger = utils.DefaultLogger
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose QUIC Connection message")
	flag.BoolVar(&quiet, "q", false, "don't print the data")
	flag.StringVar(&certPath, "cert", ".", "certificate directory")
	flag.UintVar(&repeat, "repeat", 1, "repeat time of the request")
	flag.UintVar(&interval, "interval", 1, "interval of repeat request (seconds)")
	flag.UintVar(&timeout, "timeout", 30, "idle timeout")
	flag.BoolVar(&keepAlive, "keep", false, "whether periodically send PING frames to keep the connection alive")
	flag.Parse()

	urls = flag.Args()
	logger.SetLogTimeFormat("[15:04:05.000] ")
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
		QuicConfig: &quic.Config{
			Versions: []protocol.VersionNumber{protocol.Version43}, IdleTimeout: time.Second * time.Duration(timeout), KeepAlive: keepAlive,
		},
		TLSClientConfig: comm.GetTLSConfig(certFile, keyFile),
	}
	client := &http.Client{Transport: transport}
	defer transport.Close()

	for i := uint(0); i < repeat; i++ {
		for _, addr := range urls {
			logger.Infof("GET %s", addr)

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

			logger.Infof("\n\n\n")
			time.Sleep(time.Second * time.Duration(interval))
		}
	}
}
