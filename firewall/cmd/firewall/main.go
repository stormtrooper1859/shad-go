// +build !solution

package main

import (
	"flag"
	"fmt"
	"gitlab.com/slon/shad-go/firewall/internal/firewall"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	confPath     = flag.String("conf", "C:\\\\all\\projects\\go\\shad-go\\firewall\\configs\\example.yaml", "path to config file")
	firewallAddr = flag.String("addr", ":8081", "address of firewall")
	serviceAddr  = flag.String("service-addr", "http://eu.httpbin.org/", "address of protected service")
)

func main() {
	flag.Parse()

	config, err := firewall.GetConfig(*confPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	parsedUrl, err := url.Parse(*serviceAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	reverseProxy.Transport = firewall.NewFirewall(config)

	http.Handle("/", reverseProxy)
	err = http.ListenAndServe(*firewallAddr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
