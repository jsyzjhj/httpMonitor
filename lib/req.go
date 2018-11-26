/**
请求执行类
 */
package lib

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type req struct {
	client  *http.Client
	request *http.Request
}

//临时对象池
var pipe = &sync.Pool{
	New: func() interface{} {
		return &req{
			client: &http.Client{},
		}
	},
}

//new一个req对象，并初始化默认数据
func New(method, url, data string) (rs *req) {
	rs = pipe.Get().(*req)
	body := strings.NewReader(data)
	rs.request, _ = http.NewRequest(method, url, body)
	rs.client.Timeout = 10 * time.Second //默认10s
	return
}

//为req对象设置超时时间
func (r *req) SetOvertime(overtime int) {
	r.client.Timeout = time.Duration(overtime) * time.Second
}

//为req对象设置header
func (r *req) SetHeader(header string) {
	headers := strings.Split(header, "\n") //行拆分
	for _, v := range headers {
		val := strings.Split(v, ":") //分开键和值
		if len(val) != 2 {
			continue
		}
		r.request.Header.Set(val[0], val[1])
	}
}

//执行请求，并返回状态码、内容、错误信息
func (r *req) Do() (int, string, error) {
	resp, err := r.client.Do(r.request)
	if err != nil {
		return -1, "请求执行出错", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "返回结果解析错误", err
	}
	str, err := u2s(string(body))
	if err != nil {
		str = string(body)
	}
	return resp.StatusCode, str, nil
}

//unicode转str
func u2s(str string) (string, error) {
	reg, err := regexp.Compile(`\\u[0-9a-zA-Z]{4}`)
	if err != nil {
		return "", err
	}
	strs := reg.FindAll([]byte(str), -1)
	for _, v := range strs {
		str = strings.Replace(str, string(v), word2str(string(v)), -1)
	}
	return str, nil
}

//单个编码转换
func word2str(word string) (w string) {
	w, _ = strconv.Unquote(`"` + word + `"`)
	return
}
