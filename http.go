package toolkit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
)

var (
	// 默认client，10秒超时
	defaultHttpClient = &http.Client{Timeout: time.Second * 10}
	// 日志输出
	httpLogger io.Writer = os.Stdout
)

const (
	postJsonContentType = "application/json"
	postFormContentType = "application/x-www-form-urlencoded"
)

// 日志输出。 传入nil则关闭
func HttpSetLogger(l io.Writer) {
	httpLogger = l
}

// 结构体支持
type Request struct {
	url     string
	method  string
	body    []byte
	header  map[string]string
	timeout time.Duration
}

func HttpGet(url string, params map[string]string) *Request {
	r := &Request{
		url:    url,
		method: http.MethodGet,
	}

	if len(params) > 0 {
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)

		buf.WriteString(r.url + "?")
		for k, v := range params {
			buf.WriteString(k + "=" + v + "&")
		}

		r.url = buf.String()[:buf.Len()-1]
	}

	return r
}

func HttpPostJson(url string, data any) *Request {
	r := &Request{
		url:    url,
		method: http.MethodPost,
		header: map[string]string{"Content-Type": postJsonContentType},
	}

	if data != nil {
		bs, _ := json.Marshal(data)
		r.body = bs
	}

	return r
}

func HttpPostForm(url string, data map[string]any) *Request {
	r := &Request{
		url:    url,
		method: http.MethodPost,
		header: map[string]string{"Content-Type": postFormContentType},
	}

	r.body = buildForm(data)

	return r
}

func buildForm(data map[string]any) []byte {
	if len(data) == 0 {
		return nil
	}
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	for k, v := range data {
		buf.WriteString(k + "=" + fmt.Sprintf("%v", v) + "&")
	}
	return buf.Bytes()[:buf.Len()-1]
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[key] = value
	return r
}

func (r *Request) SetTimeout(t time.Duration) *Request {
	r.timeout = t
	return r
}

func (r *Request) Do(data any) (err error) {
	start := time.Now()

	req, err := http.NewRequest(r.method, r.url, bytes.NewReader(r.body))
	if err != nil {
		r.log(start, nil, err)
		return
	}

	for k, v := range r.header {
		req.Header.Set(k, v)
	}

	client := defaultHttpClient
	if r.timeout > 0 && client.Timeout != r.timeout {
		client = &http.Client{Timeout: r.timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		r.log(start, nil, err)
		return
	}
	defer resp.Body.Close()

	// 判断状态码
	if resp.StatusCode >= http.StatusMultipleChoices {
		err = errors.New("HTTP返回状态码异常: " + resp.Status)
		return
	}

	// 解析body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log(start, nil, err)
		return
	}
	err = json.Unmarshal(b, data)
	r.log(start, b, err)
	return
}

type requestLog struct {
	Level    string          `json:"level"`
	Time     string          `json:"time"`
	Api      string          `json:"api"`
	Method   string          `json:"method"`
	Cost     string          `json:"cost"`
	Err      string          `json:"error,omitempty"`
	Body     any             `json:"body,omitempty"`
	Response json.RawMessage `json:"response,omitempty"`
}

const datetime = "2006-01-02 15:04:05"

func (r *Request) log(start time.Time, resp []byte, err error) {
	if httpLogger == nil {
		return
	}

	data := requestLog{
		Time:     start.Format(datetime),
		Api:      r.url,
		Method:   r.method,
		Cost:     strconv.Itoa(int(time.Since(start).Milliseconds())) + "ms",
		Response: resp,
	}
	if len(r.body) > 0 {
		if r.header["Content-Type"] == postJsonContentType {
			data.Body = r.body
		} else {
			data.Body = string(r.body)
		}
	}
	if err != nil {
		data.Err = err.Error()
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	buf.WriteByte('\n')

	httpLogger.Write(buf.Bytes())
}
