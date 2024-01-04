Vì ".Get() blocks until the future is ready" mình tạo 1 go routine (light weight thread) để chạy những cái bị block trong 1 thread riêng. những cái không bị block trong 1 thread riêng

-> concurrency approach 

Problem v có cách nào dùng signal không ? branch [signal](https://github.com/kingstonduy/demo-temporal/tree/signal)

```
go run start/main.go async
```

```
go run microservice/microservice.go 
```


```
go run worker1/worker1.go 
```
