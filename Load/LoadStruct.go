package Load

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func (s *SimpleServer) Address() string {
	return s.Addr
}

func (s *SimpleServer) IsAlive() bool {
	return true
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

type LoadBalancer struct {
	Port            string
	RoundRobinCount uint64
	Server          []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		Port:            port,
		Server:          servers,
		RoundRobinCount: 0,
	}
}

func (lb *LoadBalancer) GetNextServer() Server {
	// TODO: round robin selection logic
	server := lb.Server[lb.RoundRobinCount%uint64(len(lb.Server))]
	for !server.IsAlive() {
		server = lb.Server[lb.RoundRobinCount%uint64(len(lb.Server))]
		lb.RoundRobinCount++
	}
	lb.RoundRobinCount++
	return server
}

func (lb *LoadBalancer) ServerProxy(rw http.ResponseWriter, r *http.Request) {
	// TODO: proxy logic
	targetServer := lb.GetNextServer()
	fmt.Printf("Passing request to server at: %s\n", targetServer.Address())
	http.Redirect(rw, r, targetServer.Address(), http.StatusTemporaryRedirect)
}



