simulate cơ chế retry và rerun toàn bộ workflow

chạy lần đầu best case -> in ra 3 bóng đèn
chạy lần 2  khi in ra 1 bóng đèn thì bị lỗi -> fix lỗi code in ra thêm 3 bóng đèn <=> tổng cộng 4 bóng đèn
```
go run duplicate-console/started/main.go async
```


```
go run duplicate-console/microservice/microservice.go
```

```
go run duplicate-console/worker/worker.go
```