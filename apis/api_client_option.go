package apis

import (
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"time"
)

// 微信视频号API默认域名
const DefaultWXAPIHost = "https://api.weixin.qq.com"

const HttpTTL = 1 * time.Minute

var FastClient = CreateFastHttpClient()

func CreateFastHttpClient() fasthttp.Client {
	var defaultDialer = &fasthttp.TCPDialer{Concurrency: 300}

	return fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			idx := 3 // 重试三次
			for {
				idx--
				conn, err := defaultDialer.DialTimeout(addr, 10*time.Second) // tcp连接超时时间10s
				if err != fasthttp.ErrDialTimeout || idx == 0 {
					return conn, err
				}
			}
		},
	}
}

// 需初始化的参数
type (
	Options struct {
		DcsToken DcsToken // 选传参数，如果不传这个参数，则只会把 token 存在内存中，这会导致每次重启服务都要重新获取 token 和多个服务需发起多次 token 请求的问题。
		Logger   Logger   // 选传参数，不传则默认将日志直接输出在终端
	}
)

// 日志输出接口
type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// 默认日志记录器
type loggerPrint struct{}

var _ Logger = loggerPrint{}

func (loggerPrint) Info(args ...interface{}) {
	log.Println(args...)
}

func (loggerPrint) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (loggerPrint) Error(args ...interface{}) {
	log.Println(args...)
}

func (loggerPrint) Errorf(template string, args ...interface{}) {
	log.Printf(template, args...)
}
