Trong ví dụ này ta sẽ thực hiện simulate 1 service xác thực otp. 
người dùng bấm  nút -> hệ thống gen otp -> gửi otp qua mail -> người dùng input otp


Implement: tạo 1 activity dài(getOcbInfo) chạy song song vs các activity async. (best case) các activity sẽ không block nhau.

Vì sao  chạy thành công, ta thấy mình đã  khởi tạo activity dài trước tiên rồi mới tới các activity async

Problem: (worst case) sẽ như nào nếu giữa các workflow async đang chạy mà ta gọi workflow dài? [traditionalway-getting-block](https://github.com/kingstonduy/demo-temporal/tree/traditionalway-getting-block)


```
go run async/started/main.go
```

```
go run async/microservice/microservice.go 
```


```
go run async/worker/worker.go 
```