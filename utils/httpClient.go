package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func SendGetRpc(url string, headers ...map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	if len(headers) != 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("request error code:%v", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RegenTransport(pUrl string) *http.Transport{
	parseURL := mustParseURL(pUrl)
	return &http.Transport{
		MaxIdleConns:    100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout: 30 * time.Second,
		Proxy:           http.ProxyURL(parseURL),
	}
}

func mustParseURL(rawURL string) *url.URL {
	re := regexp.MustCompile(`[\r\n]`)
	rawURL = re.ReplaceAllString(rawURL, "")
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	return parsedURL
}
