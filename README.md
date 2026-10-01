# Library API

RESTful API quản lý thư viện được xây dựng bằng **Go, Gin, GORM và MySQL**.  
Project áp dụng kiến trúc phân lớp, JWT Authentication, Redis caching/rate limiting, Swagger, Docker và GitHub Actions CI.

## 1. Mục tiêu

Project được xây dựng với các mục tiêu:

- Xây dựng REST API quản lý thư viện.
- Thực hiện CRUD cho sách và tác giả.
- Quản lý người dùng và xác thực bằng JWT.
- Quản lý mượn và trả sách.
- Sử dụng MySQL làm cơ sở dữ liệu.
- Sử dụng GORM làm ORM.
- Áp dụng soft delete cho các entity phù hợp.
- Sử dụng Redis cho caching và rate limiting.
- Cung cấp Swagger để xem và kiểm thử API.
- Container hóa bằng Docker và Docker Compose.
- Tự động format, kiểm tra, test và build bằng GitHub Actions.

> Project tập trung vào REST API backend. Các phần Health Check, Graceful Shutdown và Prometheus Metrics không nằm trong phạm vi hiện tại.

---

## 2. Công nghệ sử dụng

| Công nghệ      | Vai trò                  |
| -------------- | ------------------------ |
| Go 1.24        | Ngôn ngữ/runtime backend |
| Gin            | HTTP web framework       |
| GORM           | ORM                      |
| MySQL 8        | Cơ sở dữ liệu            |
| Redis 7        | Cache và rate limiting   |
| JWT            | Authentication           |
| bcrypt         | Hash mật khẩu            |
| Swagger        | API documentation        |
| Testify        | Unit testing             |
| Docker         | Containerization         |
| Docker Compose | Chạy API + MySQL + Redis |
| GitHub Actions | Continuous Integration   |

---

## 3. Kiến trúc hệ thống

Project sử dụng kiến trúc phân lớp:

```text
Client / Postman
       |
       v
+----------------------+
|      Gin Router      |
+----------------------+
       |
       v
+----------------------+
|     Middleware       |
| JWT / Role / Rate    |
| Limit / Error        |
+----------------------+
       |
       v
+----------------------+
|       Handler        |
| HTTP Request/Response|
+----------------------+
       |
       v
+----------------------+
|       Service        |
| Business Logic       |
+----------------------+
       |
       v
+----------------------+
|     Repository       |
| Database Operations  |
+----------------------+
       |
       v
+----------------------+
|       MySQL          |
+----------------------+

Redis được sử dụng bên cạnh Service cho:
- Book list caching
- Rate limiting
```

### Request flow

```text
HTTP Request
    ↓
Router
    ↓
Middleware
    ↓
Handler
    ↓
Service
    ├── Redis Cache
    └── Repository
            ↓
          MySQL
    ↓
HTTP Response
```

---

## 4. Cấu trúc project

```text
library-api/
├── cmd/
│   └── api/
│       └── main.go
│
├── configs/
│   ├── config.go
│   ├── database.go
│   └── redis.go
│
├── internal/
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── author_handler.go
│   │   ├── book_handler.go
│   │   └── borrow_handler.go
│   │
│   ├── middlewares/
│   │   ├── auth.go
│   │   ├── error_handler.go
│   │   ├── rate_limit.go
│   │   └── role.go
│   │
│   ├── models/
│   │   ├── user.go
│   │   ├── book.go
│   │   ├── author.go
│   │   ├── borrow.go
│   │   └── return.go
│   │
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── book_repository.go
│   │   ├── author_repository.go
│   │   └── borrow_repository.go
│   │
│   ├── routes/
│   │   ├── auth_routes.go
│   │   ├── book_routes.go
│   │   ├── author_routes.go
│   │   └── borrow_routes.go
│   │
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── book_service.go
│   │   ├── author_service.go
│   │   ├── borrow_service.go
│   │   └── errors.go
│   │
│   └── validators/
│       ├── auth_request.go
│       ├── author_request.go
│       ├── book_request.go
│       ├── book_query.go
│       └── borrow_request.go
│
├── pkg/
│   ├── cache/
│   │   └── redis.go
│   ├── jwt/
│   │   └── jwt.go
│   └── response/
│       └── error.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── tests/
│   └── book_service_test.go
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 5. Cấu hình môi trường

Tạo file `.env` tại thư mục gốc:

```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=library_db

JWT_SECRET=change-this-secret-key
JWT_EXPIRE_HOURS=24

REDIS_HOST=localhost
REDIS_PORT=6379
```

Không commit `.env` lên GitHub.

File `.env.example` được sử dụng làm mẫu cấu hình.

### Khi chạy bằng Docker Compose

API chạy trong container nên phải kết nối tới service name:

```env
DB_HOST=mysql
DB_PORT=3306

REDIS_HOST=redis
REDIS_PORT=6379
```

Không dùng `localhost` cho MySQL/Redis từ bên trong container API.

---

## 6. Cài đặt và chạy local

### Yêu cầu

- Go 1.24 hoặc tương thích với project
- MySQL 8
- Redis 7
- Git

### Clone project

```powershell
git clone https://github.com/PhamDangHoan/library-api.git
cd library-api
```

### Cài dependencies

```powershell
go mod download
```

### Tạo database

Đăng nhập MySQL:

```sql
CREATE DATABASE library_db
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;
```

### Tạo `.env`

Copy `.env.example` thành `.env` và điền thông tin MySQL, JWT và Redis.

### Khởi động Redis

Nếu Redis chạy bằng Docker:

```powershell
docker start library-redis
```

Hoặc sử dụng Docker Compose như hướng dẫn ở phần Docker.

### Chạy API

```powershell
go run ./cmd/api
```

API mặc định:

```text
http://localhost:8080
```

Khi chạy thành công, terminal sẽ hiển thị trạng thái kết nối MySQL/Redis và server.

---

## 7. Chạy bằng Docker Compose

Docker Compose cung cấp:

```text
library-api
library-mysql
library-redis
```

Khởi động:

```powershell
docker compose up --build
```

Chạy background:

```powershell
docker compose up -d --build
```

Kiểm tra container:

```powershell
docker compose ps
```

Xem log:

```powershell
docker compose logs -f
```

Xem riêng API:

```powershell
docker compose logs -f api
```

Dừng:

```powershell
docker compose down
```

Dừng và xóa volumes:

```powershell
docker compose down -v
```

> `docker compose down -v` sẽ xóa dữ liệu MySQL và Redis được lưu trong Docker volumes.

API:

```text
http://localhost:8080
```

---

## 8. Authentication

API sử dụng JWT Bearer Token.

### Register

```http
POST /api/v1/auth/register
Content-Type: application/json
```

Body:

```json
{
  "name": "Normal User",
  "email": "user@example.com",
  "password": "123456"
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json
```

Body:

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

Sau khi login thành công, lấy JWT token từ response.

Các endpoint yêu cầu authentication sử dụng:

```http
Authorization: Bearer <TOKEN>
```

JWT chứa thông tin user ID, email và role.

---

## 9. API Endpoints

Base URL:

```text
http://localhost:8080
```

### Authentication

| Method | Endpoint                | Authentication | Chức năng |
| ------ | ----------------------- | -------------- | --------- |
| POST   | `/api/v1/auth/register` | Không          | Đăng ký   |
| POST   | `/api/v1/auth/login`    | Không          | Đăng nhập |

### Books

| Method | Endpoint            | Authentication | Chức năng      |
| ------ | ------------------- | -------------- | -------------- |
| GET    | `/api/v1/books`     | Không          | Danh sách sách |
| GET    | `/api/v1/books/:id` | Không          | Chi tiết sách  |
| POST   | `/api/v1/books`     | JWT            | Tạo sách       |
| PUT    | `/api/v1/books/:id` | JWT            | Cập nhật sách  |
| DELETE | `/api/v1/books/:id` | Admin          | Xóa sách       |

### Authors

| Method | Endpoint              | Authentication | Chức năng   |
| ------ | --------------------- | -------------- | ----------- |
| POST   | `/api/v1/authors`     | JWT            | Tạo tác giả |
| DELETE | `/api/v1/authors/:id` | JWT            | Xóa tác giả |

### Borrow / Return

| Method | Endpoint                     | Authentication | Chức năng        |
| ------ | ---------------------------- | -------------- | ---------------- |
| POST   | `/api/v1/borrows`            | JWT            | Mượn sách        |
| GET    | `/api/v1/borrows/my`         | JWT            | Xem sách đã mượn |
| POST   | `/api/v1/borrows/:id/return` | JWT            | Trả sách         |

> API thực tế nên được kiểm tra theo Swagger hoặc Postman Collection đi kèm project.

---

## 10. Quản lý sách

### Lấy danh sách sách

```http
GET /api/v1/books
```

Có thể sử dụng các query parameter cho pagination, search, ISBN và sorting.

Ví dụ:

```http
GET /api/v1/books?page=1&limit=10
```

Search:

```http
GET /api/v1/books?search=clean
```

Theo ISBN:

```http
GET /api/v1/books?isbn=9780132350884
```

Sorting:

```http
GET /api/v1/books?sort=title&order=asc
```

### Tạo sách

```http
POST /api/v1/books
Authorization: Bearer <TOKEN>
Content-Type: application/json
```

Body:

```json
{
  "title": "Clean Code",
  "isbn": "9780132350884",
  "description": "A handbook of agile software craftsmanship",
  "quantity": 10,
  "available": 10
}
```

Nếu tạo sách kèm tác giả:

```json
{
  "title": "Clean Code",
  "isbn": "9780132350884",
  "description": "A handbook of agile software craftsmanship",
  "quantity": 10,
  "available": 10,
  "author_ids": [1]
}
```

### Cập nhật sách

```http
PUT /api/v1/books/1
Authorization: Bearer <TOKEN>
Content-Type: application/json
```

Ví dụ:

```json
{
  "title": "Clean Code - Updated",
  "isbn": "9780132350884",
  "description": "Updated description",
  "quantity": 15,
  "available": 12
}
```

### Xóa sách

```http
DELETE /api/v1/books/1
Authorization: Bearer <ADMIN_TOKEN>
```

Book sử dụng soft delete nên bản ghi đã xóa không xuất hiện trong các truy vấn thông thường.

---

## 11. Borrow / Return

### Mượn sách

```http
POST /api/v1/borrows
Authorization: Bearer <TOKEN>
Content-Type: application/json
```

Body:

```json
{
  "book_id": 1
}
```

### Xem lịch sử mượn của user

```http
GET /api/v1/borrows/my
Authorization: Bearer <TOKEN>
```

### Trả sách

```http
POST /api/v1/borrows/1/return
Authorization: Bearer <TOKEN>
```

Khi mượn/trả, số lượng `available` của sách được cập nhật thông qua transaction.

---

## 12. Redis Caching

Redis được sử dụng làm caching layer cho danh sách sách.

Flow:

```text
GET /api/v1/books
        |
        v
   Generate Cache Key
        |
        v
      Redis
     /     \
   HIT     MISS
    |        |
    v        v
 Response   MySQL
               |
               v
             Redis
               |
               v
            Response
```

Cache key có dạng:

```text
library:books:page=1&limit=10&search=&isbn=&sort=created_at&order=desc
```

TTL mặc định của book list cache:

```text
5 minutes
```

### Cache invalidation

Khi dữ liệu sách thay đổi:

```text
POST /books
PUT /books/:id
DELETE /books/:id
```

cache pattern:

```text
library:books:*
```

được xóa để tránh trả dữ liệu cũ.

### Kiểm tra Redis

```powershell
docker exec -it library-redis redis-cli
```

Xem cache key:

```redis
KEYS "library:books:*"
```

Kiểm tra TTL:

```redis
TTL "library:books:page=1&limit=10&search=&isbn=&sort=created_at&order=desc"
```

Xóa cache khi cần:

```redis
FLUSHDB
```

> Redis chỉ đóng vai trò cache/rate limiting. MySQL vẫn là nguồn dữ liệu chính.

---

## 13. Rate Limiting

API sử dụng Redis để giới hạn request theo IP.

Giới hạn hiện tại:

```text
60 requests / minute / IP
```

Khi vượt quá giới hạn:

```http
429 Too Many Requests
```

Response:

```json
{
  "success": false,
  "message": "Too many requests. Please try again later."
}
```

API cũng trả các header liên quan đến rate limit như:

```text
X-RateLimit-Limit
X-RateLimit-Remaining
Retry-After
```

---

## 14. Swagger

Swagger documentation được sinh trong thư mục:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

Sau khi chạy API, mở:

```text
http://localhost:8080/swagger/index.html
```

Swagger cho phép xem endpoint, request schema và thực hiện thử API trực tiếp.

---

## 15. Test API bằng Postman

Project có Postman Collection:

```text
library-api.postman_collection.json
```

Import file này vào Postman:

```text
Postman
→ Import
→ library-api.postman_collection.json
```

Collection hiện có các nhóm request cho:

- Authentication
- Books
- Authors
- Borrow / Return
- Search

Collection sử dụng các biến:

```text
base_url
token
search
```

Thiết lập:

```text
base_url = http://localhost:8080
```

Sau khi login, copy JWT vào:

```text
token
```

### Trình tự test API đề xuất

#### 1. Register

```http
POST /api/v1/auth/register
```

#### 2. Login

```http
POST /api/v1/auth/login
```

Lấy JWT token.

#### 3. Tạo author

```http
POST /api/v1/authors
Authorization: Bearer <TOKEN>
```

#### 4. Tạo book

```http
POST /api/v1/books
Authorization: Bearer <TOKEN>
```

#### 5. Lấy danh sách books

```http
GET /api/v1/books
```

#### 6. Lấy book theo ID

```http
GET /api/v1/books/1
```

#### 7. Test search / pagination

```http
GET /api/v1/books?search=clean&page=1&limit=10
```

#### 8. Update book

```http
PUT /api/v1/books/1
Authorization: Bearer <TOKEN>
```

#### 9. Borrow

```http
POST /api/v1/borrows
Authorization: Bearer <TOKEN>
```

Body:

```json
{
  "book_id": 1
}
```

#### 10. Xem borrow history

```http
GET /api/v1/borrows/my
Authorization: Bearer <TOKEN>
```

#### 11. Return book

```http
POST /api/v1/borrows/1/return
Authorization: Bearer <TOKEN>
```

#### 12. Delete book

```http
DELETE /api/v1/books/1
Authorization: Bearer <ADMIN_TOKEN>
```

---

## 16. Unit Testing

Chạy toàn bộ test:

```powershell
go test ./...
```

Chạy test với verbose output:

```powershell
go test -v ./...
```

Kiểm tra coverage:

```powershell
go test ./... -cover
```

Project có unit test cho Book Service và sử dụng mock repository/cache để kiểm tra business logic mà không cần kết nối trực tiếp tới database.

---

## 17. Docker

Build image:

```powershell
docker build -t library-api .
```

Dockerfile sử dụng multi-stage build:

```text
Go Builder
    ↓
Compile application
    ↓
Small Alpine runtime image
```

Docker Compose chạy:

```text
API
├── MySQL
└── Redis
```

Khi API chạy trong Docker Compose:

```text
DB_HOST=mysql
REDIS_HOST=redis
```

---

## 18. GitHub Actions

Workflow:

```text
.github/
└── workflows/
    └── ci.yml
```

Pipeline hiện tại:

```text
Push / Pull Request
        ↓
Checkout repository
        ↓
Setup Go
        ↓
Download dependencies
        ↓
gofmt check
        ↓
go vet
        ↓
go test ./...
        ↓
go build
```

CI chạy khi có push hoặc pull request vào:

```text
main
develop
```

Mục tiêu của workflow là phát hiện sớm:

- Code chưa được format.
- Một số vấn đề được `go vet` phát hiện.
- Unit test thất bại.
- Code không build được.

---

## 19. Git workflow

Kiểm tra trạng thái:

```powershell
git status
```

Format code:

```powershell
go fmt ./...
```

Test:

```powershell
go test ./...
```

Add:

```powershell
git add .
```

Commit:

```powershell
git commit -m "feat: update library api"
```

Push:

```powershell
git push origin main
```

Sau khi push:

```text
Git
 ↓
GitHub
 ↓
GitHub Actions
 ↓
gofmt
 ↓
go vet
 ↓
go test
 ↓
go build
```

---

## 20. API testing checklist

Trước khi demo, kiểm tra:

```text
Authentication
[ ] Register
[ ] Login
[ ] JWT authentication
[ ] Role authorization

Books
[ ] Get books
[ ] Get book by ID
[ ] Search
[ ] Pagination
[ ] Sorting
[ ] Create book
[ ] Update book
[ ] Delete book

Authors
[ ] Create author
[ ] Delete author
[ ] Book ↔ Author relationship

Borrow / Return
[ ] Borrow book
[ ] View my borrows
[ ] Return book
[ ] Available quantity updated

Redis
[ ] Cache MISS
[ ] Cache HIT
[ ] TTL
[ ] Cache invalidation
[ ] Rate limiting

Documentation
[ ] Swagger
[ ] Postman Collection

Quality
[ ] gofmt
[ ] go vet
[ ] go test
[ ] go build
[ ] GitHub Actions
```

---

## 21. Repository

GitHub:

```text
https://github.com/PhamDangHoan/library-api.git
```

Postman Collection:

```text
library-api.postman_collection.json
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

---

## 22. Tác giả

**Phạm Đăng Hoàn**
**MSS:134010124044**
**Lớp:WD1306**

Project: **Library API**

Technology:

```text
Go
Gin
GORM
MySQL
Redis
JWT
Swagger
Docker
GitHub Actions
```
