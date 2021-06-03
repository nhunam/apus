# I. Hướng dẫn:
## 1. Get thư viện:
```bash
go get
```

## 2. Run docker:
```bash
cd docker-compose
docker-compose up -d
```

## 3. Run SQL script:
Để tạo bảng mẫu, location:
```
scripts/sql/schema.sql
```

## 4. Download swag:
```bash
go get -u github.com/swaggo/swag/cmd/swag
```

## 5. Run:
```bash
export ENVIRONMENT=DEV;swag init;go generate;go run .
```

## 6. Go to:
http://localhost:8081/swagger/index.html

# II. Hướng dẫn debug bằng Intellij IDEA:
## 1. Tạo cấu hình run Shell Script:
Run file scripts/pre_run.sh để init swagger resources.

![1](https://i.ibb.co/0YZT4dX/Screen-Shot-2021-03-13-at-8-52-36-PM.png)

## 2. Tạo cấu hình go build:
Run go build như hình
- bỏ trống Files
- Set Environment (DEV, PROD)
- Before launch: run Shell Script ở bước 1 và go generate (để generate Depedency Injection)

![2](https://i.ibb.co/C1448f8/Screen-Shot-2021-03-13-at-8-54-09-PM.png)

# III. Giải thích các package/folder:
## 1. api:
Là package chứa các API endpoint (controller)
## 2. cache:
Là package chứa caching function
## 3. config:
Là package chứa file load config và các file config yaml. Có suffix là môi trường.
## 4. dao:
Là packge DAO. Đảm nhiệm thao tác DB. Là layer ngay dưới package api.
## 5. database:
Là package load cấu hình database: connection, pool ...
## 6. docker-compose:
Là folder chứa file docker-compose.yml dùng để dev (postgres, redis, kafka).
## 7. docs:
Là folder khi chạy câu lệnh ```swag init``` sẽ init các file resource của Swagger.
## 8. dto:
Là package DTO, chứa DTO request và response. Chú ý API ko trả về model, phải trả về DTO.
## 9. eventsourcing:
Là package produce và consume event từ Message broker (Kafka, RabbitMQ)
## 10. i18n:
Là package chứa cấu hình và resource cho message i18n.
## 11. log:
Là package chứa cấu hình logger (zerolog)
## 12. model:
Là package chứa model entity.
## 13. operation:
Là package wrapper của DAO + cache. Sử dụng khi cần tự implement logic cache + DB (cache through)
## 14. router:
Là package chứa cấu hình router (Gin)
## 15. scripts:
Là folder chứa các script cần cho project: bash script, sql script, ... Chạy như thế nào cần mô tả trong README.
## 16. util:
Là package util của project, khai báo constant ở đây luôn.

# IV. Các chú thích khác:
## 1. Dependency Injection:
Project dùng google wire, là thư viện DI kiểu source generation. Cụ thể xem file wire.go.
Vì vậy nên các file chức năng cần khai báo kiểu OOP, dùng DI để liên kết.

## 2. Customization:
Các project có thể ko giống nhau. Nên khi có thay đổi cần lựa chọn để update lại base-service. Để cho các service clone từ base-service sử dụng.

## 3. Persistent API:
Do project dùng GORM làm ORM nên khuyến khích sử dụng GORM để thao tác với DB.

# V. Tham khảo:
## Gin graceful shutdown
https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go

## GORM Transaction
https://gorm.io/docs/transactions.html

## Tham khao
https://github.com/qiangxue/go-rest-api/blob/master/internal/healthcheck/api.go

## Implement pagination
https://medium.com/wesionary-team/implement-pagination-in-golang-using-gorm-and-gin-b4ad8e2932a6

## First-Go
https://medium.com/wesionary-team/create-your-first-rest-api-with-golang-using-gin-gorm-and-mysql-d439bcc6f987
https://github.com/SudeepTimalsina/first-go

## Redis cache
https://github.com/go-redis/cache

## Uber FX
https://github.com/uber-go/fx

## Google wire
https://github.com/google/wire

## Dependency Injection in GO with Wire
https://medium.com/wesionary-team/dependency-injection-in-go-with-wire-74f81cd222f6

https://github.com/ekhabarov/blog-code-snippets/tree/master/wire

## JWT example
https://medium.com/better-programming/hands-on-with-jwt-in-golang-8c986d1bb4c0
https://github.com/hamzawix/jwt-auth-go

## Cache decorator for Gin
https://github.com/gin-contrib/cache/blob/master/cache.go

## Gin JWT middleware
https://github.com/appleboy/gin-jwt

## Gin contributor cache
https://github.com/gin-contrib/cache

## Kafka-go
https://yusufs.medium.com/getting-started-with-kafka-in-golang-14ccab5fa26
https://github.com/segmentio/kafka-go

## Gin Swagger middleware
https://github.com/swaggo/gin-swagger
https://github.com/swaggo/swag/tree/master/example/celler
https://golangexample.com/automatically-generate-restful-api-documentation-with-swagger-2-0-for-go/

## Dockerfile
https://github.com/segmentio/kafka-go/blob/master/examples/consumer-logger/Dockerfile
https://codefresh.io/docs/docs/learn-by-example/golang/golang-hello-world/
https://tutorialedge.net/golang/go-docker-tutorial/
https://blog.golang.org/docker

## Zerolog
https://medium.com/swlh/stop-using-prints-and-embrace-zerolog-2c4dd8e8816a
https://github.com/rs/zerolog

## Resty
https://github.com/go-resty/resty

## Gin CORS
https://github.com/gin-contrib/cors

## RabbitMQ
https://x-team.com/blog/set-up-rabbitmq-with-docker-compose/
https://github.com/rabbitmq/rabbitmq-tutorials/tree/master/go
Cách áp dụng group và partition vào RabbitMQ (mặc định ko support):
https://docs.spring.io/spring-cloud-stream-binder-rabbit/docs/current/reference/html/spring-cloud-stream-binder-rabbit.html#_partitioning_with_the_rabbitmq_binder

