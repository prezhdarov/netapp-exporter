package ontapi

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/prezhdarov/prometheus-exporter/collector"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	naUser     = flag.String("ontap.username", "", "Username to login to ONTAP Cluster")
	naPassword = flag.String("ontap.password", "", "Password for the user above")
	naCluster  = flag.String("ontap.cluster", "", "ONTAP Cluster REST API address in host:port format.")
	naSchema   = flag.String("ontap.schema", "https", "Use HTTP or HTTPS")
	naSSL      = flag.Bool("ontap.ssl", false, "Verify ONTAP SSL or trust")
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

func Load(logger log.Logger) {

	level.Info(logger).Log("msg", "Loading ONTAP REST API")

}

func (na *ONTAP) Login(target string) (map[string]interface{}, error) {

	var err error

	loginData := make(map[string]interface{}, 0)

	if target == "" {

		target = *naCluster

	}

	loginData["target"] = target
	loginData["headers"] = map[string]string{"Accept": "application/json"}

	loginData["nodes"], err = na.listNodes(loginData)
	if err != nil {

		return nil, err

	}

	return loginData, nil

}

// NetApp ONTAP REST API doesn't believe in sessions and suchlike therefore there isn't much that can be done to logout.... doing what we can, aren't we :D
func (na *ONTAP) Logout(loginData map[string]interface{}) error {

	return nil

}

func (na *ONTAP) Get(loginData, extraConfig map[string]interface{}) (interface{}, error) {

	url := fmt.Sprintf("%s://%s%s", *naSchema, loginData["target"], extraConfig["api"])

	_, _, body, err := request("GET", url, loginData["headers"].(map[string]string), []string{})
	if err != nil {
		return nil, err
	}

	return &body, nil
}

//request is where the http magic happens
func request(method, url string, headers map[string]string, responseHeaders []string) (int, map[string]string, []byte, error) {
	resHeaders := map[string]string{}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !*naSSL}},

		Timeout: 0,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, nil, nil, err
	}

	req.SetBasicAuth(*naUser, *naPassword)

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

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, resHeaders, body, nil
}
