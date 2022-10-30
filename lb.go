package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Backend struct {
	Ip        string
	Port      string
	NumReq    int
	IsHealthy bool
	mu        sync.RWMutex
}

type Lb struct {
	Backends []*Backend
	Strategy RoundRobin
}

type IncomingReq struct {
	sourceConn net.Conn
	timestamp  int64
}

type RoundRobin struct {
	Backends []*Backend
	Index    int
}

func InitRR(b []*Backend) {

	strategy = &RoundRobin{
		Backends: b,
		Index:    0,
	}

	lb.Strategy = *strategy
}

func InitLb() {
	backends := []*Backend{
		&Backend{Ip: "localhost", Port: "8080", NumReq: 0, IsHealthy: true},
		&Backend{Ip: "localhost", Port: "8081", NumReq: 0, IsHealthy: true},
		&Backend{Ip: "localhost", Port: "8082", NumReq: 0, IsHealthy: true},
		&Backend{Ip: "localhost", Port: "8083", NumReq: 0, IsHealthy: true},
	}

	lb = &Lb{
		Backends: backends,
	}

	InitRR(backends)
}

func (b *Backend) SetHealthStatus(status bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.IsHealthy = status
}

func (b *Backend) GetHealthStatus() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.IsHealthy
}

func (strategy *RoundRobin) GetBackend() *Backend {
	strategy.Index = (strategy.Index + 1) % len(strategy.Backends)
	return strategy.Backends[strategy.Index]
}

func (lb *Lb) Run() {
	lb_server, err := net.Listen("tcp", ":8000")

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Load Balancer is listening on port 8000")
	defer lb_server.Close()
	for {
		source_connection, err := lb_server.Accept()
		if err != nil {
			fmt.Println("Error connecting to the client", '\n')
		}

		go lb.Forward(IncomingReq{
			sourceConn: source_connection,
			timestamp:  time.Time.Unix(time.Now()),
		})
	}
}

func (lb *Lb) Forward(req IncomingReq) {
	backend := lb.Strategy.GetBackend()

	backendConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", backend.Ip, backend.Port))
	if err != nil {
		fmt.Printf("Error connecting to backend server at port :: %s\n", backend.Port)
		healthStatus := backend.GetHealthStatus()
		if healthStatus == true {
			backend.SetHealthStatus(false)
		}
		for _, v := range lb.Backends {
			fmt.Println("---------------------------")
			fmt.Println(v.GetHealthStatus())
			fmt.Println("---------------------------")
		}
		req.sourceConn.Write([]byte("Server is down"))
		req.sourceConn.Close()
		return
	}

	fmt.Printf("Request routed to :: %s:%s\n", backend.Ip, backend.Port)

	if backendConn != nil && backend.GetHealthStatus() != true {
		backend.SetHealthStatus(true)
	}

	backend.NumReq++

	for _, v := range lb.Backends {
		fmt.Println("---------------------------")
		fmt.Println(v)
		fmt.Println("---------------------------")
	}
	go io.Copy(backendConn, req.sourceConn)
	go io.Copy(req.sourceConn, backendConn)
}

var lb *Lb
var strategy *RoundRobin

func main() {
	InitLb()
	lb.Run()
}
