package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type proxy struct {
	table *table
}

func NewProxy() *proxy {
	return &proxy{
		table: newTable(),
	}
}

//管理逻辑
func (p *proxy) admin(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	if data != "" {
		//重置路由表
		p.table.DelAll()
		arr := strings.Split(data, "\r\n")
		for _, item := range arr {
			d := strings.Split(item, "=>")
			if len(d) < 2 || d[0] == "admin" {
				continue
			}
			_, err := url.Parse("http://" + d[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(d[0], d[1])
			p.table.Set(d[0], d[1])
		}
	}
	var newData string
	for service, newUrl := range p.table.GetAll() {
		if service == "" {
			continue
		}
		newData = newData + service + "=>" + newUrl + "\n"
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<html><head><title>开发联调神器</title></head><form method=\"POST\"><center>开发联调神器<br><textarea placeholder=\"test=>192.168.8.8:8080\r\n效果：\r\n" + GetURL(r) + "/test/product/list => http://192.168.8.8:8080/product/list\" autofocus name=\"data\" rows=\"30\" cols=\"100\">" + newData + "</textarea><br><input type=\"submit\" value=\"提交\"></center></form><html>"))
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	u := strings.Split(r.URL.Path, "/")
	service := u[1]
	//管理后台
	if service == "admin" {
		p.admin(w, r)
		return
	}
	s := p.table.Get(service)
	if s == "" {
		w.Write([]byte("没有设置信息，请检查配置：" + GetURL(r) + "/admin"))
		return
	}
	target, err := url.Parse("http://" + p.table.Get(service))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/"+service)
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Host = target.Host
		req.Header.Set("X-Real-Ip", req.RemoteAddr)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

//获取url中的第一个参数
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func GetURL(r *http.Request) (Url string) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host}, "")
}
