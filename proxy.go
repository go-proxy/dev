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
	//重置路由表
	p.table.DelAll()
	data := r.FormValue("data")
	arr := strings.Split(data, "\n")
	var newArr []string
	for _, item := range arr {
		d := strings.Split(item, "=>")
		if len(d) < 2 {
			continue
		}
		_, err := url.Parse(d[1])
		if err != nil {
			continue
		}
		newArr = append(newArr, item)
		fmt.Println(d[0], d[1])
		p.table.Set(d[0], d[1])
	}
	newData := strings.Join(newArr, "\n")
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<form><center><textarea autofocus name=\"data\" rows=\"30\" cols=\"100\">" + newData + "</textarea><br><input type=\"submit\" value=\"提交\"></center></form>"))
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
	target, err := url.Parse(p.table.Get(service))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	r.URL.Path = strings.TrimLeft(r.URL.Path, "/"+service)
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
