package zhima

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

func send(httpClient *http.Client, url string, body io.Reader, opt *Options) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return sendWithContext(ctx, httpClient, url, body, opt)
}

// Sending an HTTP request and accepting context.
func sendWithContext(ctx context.Context, httpClient *http.Client, url string, body io.Reader, opt *Options) (*http.Response, error) {
	v, _ := query.Values(opt)

	// fmt.Print(v.Encode()) will output: "city=0&mr=1&pb=4&pro=0&yys=0"
	APIEndpoint := fmt.Sprintf("%s&%s", url, v.Encode())
	fmt.Println(APIEndpoint)
	// Change NewRequest to NewRequestWithContext and pass context it
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIEndpoint, body)
	if err != nil {
		return nil, err
	}
	// http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
