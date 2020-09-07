package zhima

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	URL          = "http://http.tiqu.alicdns.com/getip3?num=1&type=2&ts=1&ys=1&cs=1&lb=1&sb=0"
	SpeedTestURL = "http://proxies.site-digger.com/proxy-detect/"
)

type Options struct {
	Pro  int `url:"pro"`  // 省份，默认全国
	City int `url:"city"` // 城市，默认全国
	YYS  int `url:"yys"`  // 0:不限 100026:联通 100017:电信
	MR   int `url:"mr"`   // 去重选择（1:360天去重 2:单日去重 3:不去重）
	PB   int `url:"pb"`   // 端口位数（4:4位端口 5:5位端口）
	Port int `url:"port"` // IP协议 1:HTTP 2:SOCK5 11:HTTPS
	Time int `url:"time"` // 稳定时长
}

type Proxy struct {
	IP         string `json:"ip"`          // 隧道 ip (代理ip)
	Port       int    `json:"port"`        // 代理端口
	ExpireTime string `json:"expire_time"` // 过期时间
	City       string `json:"city"`        // 城市
	ISP        string `json:"isp"`         // 运营商
	Outip      string `json:"outip"`       // 隧道 ip 的出口 ip
}

type Response struct {
	Code    int     `json:"code"`
	Data    []Proxy `json:"data"`
	Msg     string  `json:"msg"`
	Success bool    `json:"success"`
}

func GetProxy(opt Options) (proxy Proxy, err error) {
	resp, err := send(bytes.NewBuffer(nil), &opt)
	if err != nil {
		log.Println(err)
		return
	}

	// read all response body
	// data, _ := ioutil.ReadAll(resp.Body)
	// log.Println(data)
	// fmt.Printf("%s\n", data)

	// Close response body
	defer resp.Body.Close()
	proxy, err = decoder(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return proxy, nil
}

func TestProxy(proxy Proxy) (speed, status int) {
	// 解析代理地址
	proxyAddr, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	// 设置网络传输.
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxyAddr),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	// 创建连接客户端
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	begin := time.Now() // 判断代理访问时间
	// 使用代理IP访问测试地址
	res, err := httpClient.Get(SpeedTestURL)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	defer res.Body.Close()

	speed = int(time.Since(begin).Nanoseconds() / 1000 / 1000) // ms
	// 判断是否成功访问，如果成功访问 StatusCode 应该为 200
	if res.StatusCode != http.StatusOK {
		fmt.Println(err)
		return
	}
	return speed, res.StatusCode
}

func decoder(body io.Reader) (proxy Proxy, err error) {
	var resBody Response
	err = json.NewDecoder(body).Decode(&resBody)
	if err != nil {
		return
	}

	if !resBody.Success {
		// nolinter
		return proxy, errors.New(resBody.Msg)
	}

	proxy = resBody.Data[0]
	return proxy, nil
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
