package zhima

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-querystring/query"
)

const URL = "http://http.tiqu.alicdns.com/getip3?num=1&type=2&port=1&time=1&ts=1&ys=0&cs=0&lb=1&sb=0"

type Options struct {
	Pro  int `url:"pro"`
	City int `url:"city"`
	YYS  int `url:"yys"`
	MR   int `url:"mr"`
	PB   int `url:"pb"`
}

type RespData struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	ExpireTime string `json:"expire_time"`
	Outip      string `json:"outip"`
}

type Response struct {
	Code    int        `json:"code"`
	Data    []RespData `json:"data"`
	Msg     string     `json:"msg"`
	Success bool       `json:"success"`
}

func GetIP(opt Options) (result string, err error) {
	resp, err := send(bytes.NewBuffer(nil), &opt)
	if err != nil {
		log.Println(err)
		return
	}

	// Close response body
	defer resp.Body.Close()
	result, err = decoder(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return result, nil
}

func decoder(body io.Reader) (result string, err error) {
	var resBody Response
	err = json.NewDecoder(body).Decode(&resBody)
	if err != nil {
		return
	}

	Data := resBody.Data[0]

	result = fmt.Sprintf("http://%s:%d", Data.IP, Data.Port)

	return result, nil
}

func send(body io.Reader, opt *Options) (*http.Response, error) {
	return sendWithContext(context.Background(), body, opt)
}

// Sending an HTTP request and accepting context.
func sendWithContext(ctx context.Context, body io.Reader, opt *Options) (*http.Response, error) {
	v, _ := query.Values(opt)

	// fmt.Print(v.Encode()) will output: "city=0&mr=1&pb=4&pro=0&yys=0"
	APIEndpoint := fmt.Sprintf("%s&%s", URL, v.Encode())
	fmt.Println(APIEndpoint)
	// Change NewRequest to NewRequestWithContext and pass context it
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIEndpoint, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
