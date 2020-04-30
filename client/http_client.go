package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/itea-tgl/itea-go/constant"
	"github.com/itea-tgl/itea-go/ilog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GET_REQUEST_TIMEOUT 	= 3
	POST_REQUEST_TIMEOUT 	= 5
	DOWNLOAD_REQUEST_TIMEOUT 	= 10
)

type HttpClient struct {
	Ctx 		context.Context
	debug 		bool
	Proxy   	string	`value:"client.http.proxy"`
	SkipHttps 	bool	`value:"client.http.skip_https"`
}

func (c *HttpClient) Construct() {
	c.debug = c.Ctx.Value(constant.DEBUG).(bool)
}

type RequestBody struct {
	Method string
	Uri string
	Header map[string]string
	Params map[string]string
	SParam string
	Host string
	Timeout int
}

type ResponseBody struct {
	Code int
	Body []byte
}

func (c *HttpClient) Get(r *RequestBody) (res *ResponseBody, err error) {
	if c.debug {
		start := time.Now()
		defer func() {
			ilog.Info(fmt.Sprintf("【GET请求】耗时：%s, 地址[%s]", time.Since(start).String(), r.Uri))
		}()
	}

	if r.Method == "" {
		r.Method = "GET"
	}
	if r.Timeout <= 0 {
		r.Timeout = GET_REQUEST_TIMEOUT
	}
	request, err := http.NewRequest(r.Method, r.Uri, nil)
	if err != nil {
		return
	}
	if r.Params != nil {
		q := request.URL.Query()
		for k, v := range r.Params {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}
	if !strings.EqualFold(r.Host, "") {
		request.Host = r.Host
	}
	for k, v := range r.Header {
		request.Header.Set(k, v)
	}
	response, err := c.client(r.Timeout).Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	res = &ResponseBody{
		Code: response.StatusCode,
		Body: body,
	}
	return
}

func (c *HttpClient) Post(r *RequestBody) (res *ResponseBody, err error) {
	if c.debug {
		start := time.Now()
		defer func() {
			ilog.Info(fmt.Sprintf("【POST请求】耗时：%s, 地址[%s]", time.Since(start).String(), r.Uri))
		}()
	}

	if r.Method == "" {
		r.Method = "POST"
	}
	if r.Timeout <= 0 {
		r.Timeout = POST_REQUEST_TIMEOUT
	}

	postParams := url.Values{}
	for k, v := range r.Params {
		postParams.Set(k, v)
	}

	request, err := http.NewRequest(r.Method, r.Uri, strings.NewReader(postParams.Encode()))
	if err != nil {
		return
	}
	if !strings.EqualFold(r.Host, "") {
		request.Host = r.Host
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range r.Header {
		request.Header.Set(k, v)
	}

	response, err := c.client(r.Timeout).Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	res = &ResponseBody{
		Code: response.StatusCode,
		Body: body,
	}
	return
}

//func (c *HttpClient) Get(u string, h map[string]string, host string, timeout int) ([]byte, error) {
//	var body []byte
//
//	if c.debug {
//		start := time.Now()
//		defer func() {
//			ilog.Info("【GET请求】耗时：", time.Since(start), ", 地址[", u, "], 返回[", string(body),"]")
//		}()
//	}
//
//	if timeout <= 0 {
//		timeout = GET_REQUEST_TIMEOUT
//	}
//
//	body, err := c.doGet(u, h, host, timeout)
//	return body, err
//}

//func (c *HttpClient) Post(u string, p map[string]string, h map[string]string, host string, timeout int) ([]byte, error) {
//	var body []byte
//
//	if c.debug {
//		start := time.Now()
//		defer func() {
//			ilog.Info("【POST请求】耗时：", time.Since(start), ", 地址[", u, "], 参数[", p, "], 返回[", string(body),"]")
//		}()
//	}
//
//	if timeout <= 0 {
//		timeout = POST_REQUEST_TIMEOUT
//	}
//
//	postParams := url.Values{}
//	for k, v := range p {
//		postParams.Set(k, v)
//	}
//
//	body, err := c.doPost(u, h, host, timeout, "application/x-www-form-urlencoded", strings.NewReader(postParams.Encode()))
//	return body, err
//}

//func (c *HttpClient) PostJson(u string, p string, h map[string]string, host string, timeout int) ([]byte, error) {
//	var body []byte
//
//	if c.debug {
//		start := time.Now()
//		defer func() {
//			ilog.Info("【POST请求】耗时：", time.Since(start), ", 地址[", u, "], 参数[", p, "], 返回[", string(body),"]")
//		}()
//	}
//
//	if timeout <= 0 {
//		timeout = POST_REQUEST_TIMEOUT
//	}
//
//	body, err := c.doPost(u, h, host, timeout, "application/json;charset=UTF-8", strings.NewReader(p))
//	return body, err
//}

//func (c *HttpClient) PostFile(u string, file string, filekey string, p map[string]string, h map[string]string, host string, timeout int, skipHttps bool) ([]byte, error) {
//	var body []byte
//
//	if c.debug {
//		start := time.Now()
//		defer func() {
//			ilog.Info("【POST FILE请求】耗时：", time.Since(start), ", 地址[", u, "], 参数[", p, "], 文件[", file ,"], 返回[", string(body),"]")
//		}()
//	}
//
//	if timeout <= 0 {
//		timeout = POST_REQUEST_TIMEOUT
//	}
//
//	//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
//	bodyBuf := &bytes.Buffer{}
//	bodyWriter := multipart.NewWriter(bodyBuf)
//	fileWriter, err := bodyWriter.CreateFormFile(filekey, filepath.Base(file))
//	if err != nil {
//		return nil, err
//	}
//
//	//打开文件
//	f, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//
//	//把文件流写入到缓冲区里去
//	_, err = io.Copy(fileWriter, f)
//	if err != nil {
//		return nil, err
//	}
//
//	for k, v := range p {
//		bodyWriter.WriteField(k, v)
//	}
//
//	contentType := bodyWriter.FormDataContentType()
//
//	bodyWriter.Close()
//
//	body, err = c.doPost(u, h, host, timeout, contentType, ioutil.NopCloser(bodyBuf))
//	return body, err
//}

//func (c *HttpClient) Download(u string, h map[string]string, host string, timeout int, skipHttps bool) ([]byte, error) {
//	if c.debug {
//		start := time.Now()
//		defer func() {
//			ilog.Info("【DOWNLOAD请求】耗时：", time.Since(start), ", 地址[", u, "]")
//		}()
//	}
//
//	if timeout <= 0 {
//		timeout = DOWNLOAD_REQUEST_TIMEOUT
//	}
//
//	return c.doGet(u, h, host, timeout)
//}

func (c *HttpClient) client(timeout int) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	tr := &http.Transport{}

	if c.SkipHttps {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if c.Proxy != "" {
		p, _ := url.Parse(c.Proxy)
		tr.Proxy = http.ProxyURL(p)
	}

	client.Transport = tr

	return client
}

//func (c *HttpClient) doGet(u string, h map[string]string, host string, timeout int) ([]byte, error) {
//	client := c.client(timeout)
//
//	req, err := http.NewRequest("GET", u, strings.NewReader(""))
//	if err != nil {
//		return nil, err
//	}
//
//	if !strings.EqualFold(host, "") {
//		req.Host = host
//	}
//
//	for k, v := range h {
//		req.Header.Set(k, v)
//	}
//
//	resp, err := client.Do(req)
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}

//func (c *HttpClient) doPost(u string, h map[string]string, host string, timeout int, contentType string, reader io.Reader) ([]byte, error) {
//	client := c.client(timeout)
//
//	req, err := http.NewRequest("POST", u, reader)
//	if err != nil {
//		return nil, err
//	}
//
//	if !strings.EqualFold(host, "") {
//		req.Host = host
//	}
//
//	req.Header.Set("Content-Type", contentType)
//
//	for k, v := range h {
//		req.Header.Set(k, v)
//	}
//
//	resp, err := client.Do(req)
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return body, nil
//}