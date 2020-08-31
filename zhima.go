package zhima

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const URL = "http://webapi.http.zhimacangku.com/getip?num=1&type=1&pro=110000&city=110105&yys=0&port=1&time=2&ts=1&ys=0&cs=0&lb=1&sb=0&pb=4&mr=3&regions=110000"

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

func GetIP() (result string) {
	resp, err := send(bytes.NewBuffer(nil))
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

	return result
}

func decoder(body io.Reader) (result string, err error) {
	var resBody Response
	err = json.NewDecoder(body).Decode(&resBody)
	if err != nil {
		log.Println(err)
		return "", err
	}

	Data := resBody.Data[0]

	result = fmt.Sprintf("http://%s:%d", Data.IP, Data.Port)

	return result, nil
}

func send(body io.Reader) (*http.Response, error) {
	return sendWithContext(context.Background(), body)
}

// Sending an HTTP request and accepting context.
func sendWithContext(ctx context.Context, body io.Reader) (*http.Response, error) {
	// Change NewRequest to NewRequestWithContext and pass context it
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
