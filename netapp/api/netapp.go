package ontapi

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"time"

	"net/http"

	"github.com/prezhdarov/prometheus-exporter/collector"
)

var (
	naUsername = flag.String("ontap.username", "", "Username to login to ONTAP Cluster")
	naPassword = flag.String("ontap.password", "", "Password for the user above")
	naCluster  = flag.String("ontap.cluster", "", "ONTAP Cluster REST API address in host:port format.")
	naSchema   = flag.String("ontap.schema", "https", "Use HTTP or HTTPS")
	naSSL      = flag.Bool("ontap.ssl", false, "Verify ONTAP SSL or trust")
	naTimeout  = flag.Int("spo.timeout", 10, "Time in seconds to wait for NetApp Clusterr to reply")
)

type ONTAP struct {
	//logger  log.Logger
}

func init() {

	collector.RegisterAPI(NewAPI())

}

func NewAPI() *ONTAP {

	return &ONTAP{}
}

func Load(logger *slog.Logger) {

	logger.Info("msg", "Loading ONTAP REST API", nil)

}

func (na *ONTAP) Login(target string, logger *slog.Logger) (map[string]interface{}, error) {

	loginData := make(map[string]interface{}, 0)

	if target == "" {

		target = *naCluster

	}

	loginData["target"] = target
	loginData["headers"] = map[string]string{"Accept": "application/json"}

	return loginData, nil

}

// NetApp ONTAP REST API doesn't believe in sessions and suchlike therefore there isn't much that can be done to logout.... doing what we can, aren't we :D
func (na *ONTAP) Logout(loginData map[string]interface{}, logger *slog.Logger) error {

	return nil

}

func (na *ONTAP) Get(loginData, extraConfig map[string]interface{}, logger *slog.Logger) (interface{}, error) {

	url := fmt.Sprintf("%s://%s%s", *naSchema, loginData["target"], extraConfig["api"])

	_, _, body, err := request("GET", url, loginData["headers"].(map[string]string), []string{})
	if err != nil {
		return nil, err
	}

	return &body, nil
}

// request is where the http magic happens
func request(method, url string, headers map[string]string, responseHeaders []string) (int, map[string]string, []byte, error) {

	resHeaders := map[string]string{}

	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: !*naSSL}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(*naTimeout) * time.Second,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, nil, nil, err
	}

	req.SetBasicAuth(*naUsername, *naPassword)

	for header := range headers {
		req.Header.Add(header, headers[header])
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	for responseHeader := range responseHeaders {

		resHeaders[responseHeaders[responseHeader]] = resp.Header.Get(responseHeaders[responseHeader])
	}
	//headers := resp.Header.Get()
	//fmt.Println(resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, resHeaders, body, nil
}
