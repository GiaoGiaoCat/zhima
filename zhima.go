package zhima

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

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
	httpClient := http.DefaultClient
	resp, err := send(httpClient, URL, bytes.NewBuffer(nil), &opt)
	if err != nil {
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
		return
	}

	return proxy, nil
}

func TestProxy(proxy Proxy) (speed, status int, err error) {
	// 解析代理地址
	proxyAddr, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		return
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

	begin := time.Now()                                                          // 判断代理访问时间
	res, err := send(httpClient, SpeedTestURL, bytes.NewBuffer(nil), &Options{}) // 使用代理IP访问测试地址
	if err != nil {
		return
	}
	defer res.Body.Close()

	speed = int(time.Since(begin).Nanoseconds() / 1000 / 1000) // 单位 ms

	if res.StatusCode != http.StatusOK {
		return
	}

	return speed, res.StatusCode, nil
}

func decoder(body io.Reader) (proxy Proxy, err error) {
	var resBody Response
	err = json.NewDecoder(body).Decode(&resBody)
	if err != nil {
		return
	}

	if !resBody.Success {
		return proxy, errors.New(resBody.Msg)
	}

	proxy = resBody.Data[0]
	return proxy, nil
}
