# dev-proxy
开发调试神器

# 设置示例
```
test=>http://127.0.0.1:8080

假如服务域名为 dev.com 映射如下
dev.com/test/product/list => http://192.168.8.8:8080/product/list

注意:dev.com机器能访问 IP 192.168.8.8
```

# 使用方法
```
package main

import (
	"flag"
	"github.com/go-proxy/dev"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("port")
	if port == "" {
		p := flag.String("port", "8888", "port default 8888")
		port = *p
	}
	dev := proxy.NewProxy()
	log.Println("start port :" + port)
	err := http.ListenAndServe(":"+port, dev)
	if err != nil {
		log.Fatal(err)
	}
}
```