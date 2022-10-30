# Load Balancer

This is a toy load balancer, written in Golang. It performs simple round robin over a set of servers.

To run this:

```
go run lb.go
```

And the load balancer should be up and running at `localhost:8000`

To see it in action,

Run netcat commands to make a request to the load balancer

`nc 127.0.0.1 8000`