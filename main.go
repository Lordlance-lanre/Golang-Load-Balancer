package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Lordlance-lanre/Golang-Load-Balancer/Load"
)

func newSimpleServer(addr string) *Load.SimpleServer {
	serverUrl, err := url.Parse(addr)
	handleError(err)
	return &Load.SimpleServer{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		panic(err)
	}
}

func main() {
	fmt.Println("Hello World")

	servers := []Load.Server{
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.youtube.com"),
		newSimpleServer("https://www.reddit.com"),
	}

	lb := Load.NewLoadBalancer(":3000", servers)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServerProxy(rw, r)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Starting Load Balancer at port %s\n", lb.Port)
	http.ListenAndServe(lb.Port, nil)
}






 // code  to see the load balancer alternating between servers on every refresh:

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"

// 	"github.com/Lordlance-lanre/Golang-Load-Balancer/Load"
// )

// func newSimpleServer(addr string) *Load.SimpleServer {
// 	serverUrl, err := url.Parse(addr)
// 	handleError(err)

// 	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
// 	// Rewrite the Host header so backends receive the correct host
// 	originalDirector := proxy.Director
// 	proxy.Director = func(req *http.Request) {
// 		originalDirector(req)
// 		req.Host = serverUrl.Host
// 	}

// 	return &Load.SimpleServer{
// 		Addr:  addr,
// 		Proxy: proxy,
// 	}
// }

// func handleError(err error) {
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err.Error())
// 		panic(err)
// 	}
// }

// func startMockServer(port string, name string) {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// Ignore favicon to keep the round-robin count clean
// 		if r.URL.Path == "/favicon.ico" {
// 			return
// 		}
// 		fmt.Fprintf(w, "Response served from: %s (Port %s)\n", name, port)
// 	})
// 	go http.ListenAndServe(port, mux)
// }

// func main() {
// 	// Start 3 mock backend servers
// 	startMockServer(":8081", "Backend Server 1")
// 	startMockServer(":8082", "Backend Server 2")
// 	startMockServer(":8083", "Backend Server 3")

// 	servers := []Load.Server{
// 		newSimpleServer("http://localhost:8081"),
// 		newSimpleServer("http://localhost:8082"),
// 		newSimpleServer("http://localhost:8083"),
// 	}

// 	lb := Load.NewLoadBalancer(":3000", servers)

// 	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
// 		// Skip favicon requests so each browser refresh advances by exactly 1
// 		if r.URL.Path == "/favicon.ico" {
// 			return
// 		}
// 		lb.ServerProxy(rw, r)
// 	})

// 	fmt.Printf("Starting Load Balancer at port %s\n", lb.Port)
// 	http.ListenAndServe(lb.Port, nil)
// }
