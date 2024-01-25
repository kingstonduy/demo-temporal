```
go run saga/t24-service/controller/controller.go 
go run saga/napas-service/controller/controller.go
go run saga/limitation-manage-service/controller/controller.go 
go run saga/money-transfer-service/started/main.go
go run saga/money-transfer-service/worker/main.go
```

Basic case, khi execute activity gọi đến 1 service khác thất bại, sẽ retry 1 số lần nhất định, nếu sau khi retry vẫn ko được thì tự động assume là service đó đã rollback rồi (hoặc chưa có gì thay đổi ở database service đó), chỉ compensate các activity trước đó thôi. 