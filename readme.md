# Load Balancer

This is a toy load balancer, written in Golang. It performs simple round robin over a set of servers.

To run this:

```
go run lb.go
```

And the load balancer should be up and running at `localhost:8000`

To see it in action,

1. Spin four local servers on ports : 8080, 8081, 8082, 8083 (I used node to spin the servers)
2. Run netcat commands to make a request to the load balancer

`nc 127.0.0.1 8000`

Check the logs in the terminal to see the output.

> It also performs health checks on every 1 minutes and show the status of all the backend servers.
