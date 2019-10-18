package httplib

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var HttpClient *http.Client

func init() {
	HttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: func(netr, addr string) (net.Conn, error) {
				conn, e := net.DialTimeout(netr, addr, time.Second*5)
				if e != nil {
					return nil, e
				}
				conn.SetDeadline(time.Now().Add(time.Second * 5))
				return conn, nil
			},
			MaxIdleConns:          100,
			ResponseHeaderTimeout: time.Second * 5,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := HttpClient.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
