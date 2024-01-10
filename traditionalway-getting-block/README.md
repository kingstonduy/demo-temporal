implement: Cách làm cũ sẽ ko chạy được vì expect là trong lúc chờ input nó sẽ chạy các thằng acitivity khác

Fix: tạo 1 thread cho activity cần input và 1 thread cho các non-blocking activity branch [goroutines-activity](https://github.com/kingstonduy/demo-temporal/tree/goroutines-activities)


```
go run traditionalway-getting-block/started/main.go
```

```
go run traditionalway-getting-block/microservice/microservice.go 
```


```
go run traditionalway-getting-block/worker/worker.go 
```
v