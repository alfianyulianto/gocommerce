# GoCommerce

## ✨ Fitur Utama

- **User Management**: Registrasi, autentikasi, dan manajemen profil user
- **Product Catalog**: CRUD operasi untuk produk dengan pencarian full-text menggunakan Elasticsearch
- **Orders**: Sistem pemesanan
- **Full-Text Search**: Pencarian cepat untuk users, products, dan orders menggunakan Elasticsearch
- **Event Streaming**: Real-time event processing menggunakan Apache Kafka
- **File Upload**: Upload dan manajemen file (gambar produk, dll)
- **RESTful API**: API yang konsisten dan mudah digunakan

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.25.0
- **Web Framework**: Fiber v3.2.0 (HTTP server yang cepat)
- **Database**: 
  - MySQL 9 (primary database)
  - GORM (ORM untuk database operations)
- **Search Engine**: Elasticsearch 8.14.0
- **Message Queue**: Apache Kafka (latest)

### Tools & Libraries
- **Validation**: go-playground/validator/v10
- **Configuration**: Spf13/Viper
- **Logging**: Sirupsen/Logrus
- **UUID**: Google/UUID
- **OpenTelemetry**: untuk observability

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database Migration**: SQL migrations

## 📦 Prasyarat

Sebelum memulai, pastikan Anda telah menginstal:

- **Go**: v1.25.0 atau lebih baru
- **Docker & Docker Compose**: Untuk menjalankan services (MySQL, Elasticsearch, Kafka, Kibana)
- **Git**: Untuk clone repository

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone https://github.com/alfianyulianto/gocommerce.git
cd gocommerce
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Environment Variables

Copy file `.env.example` ke `.env`:
```bash
cp .env.example .env
```

Contoh konfigurasi default untuk development:
```env
APP_NAME=GoCommerce
APP_ENV=development
APP_PORT=3000
APP_BASE_URL=http://localhost:3000

DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_NAME=gocommerce
DATABASE_USER=root
DATABASE_PASSWORD=root

POOL_MAX_IDLE_CONNS=10
POOL_MAX_OPEN_CONNS=100
POOL_CONN_IDLE_TIME=600
POOL_CONN_LIFE_TIME=3600

ELASTICSEARCH_ADDRESSES=http://localhost:9201
ELASTICSEARCH_USER_INDEX=users
ELASTICSEARCH_PRODUCT_INDEX=products
ELASTICSEARCH_ORDER_INDEX=orders

KAFKA_GROUP_ID=gocommerce-group
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=order-events
```

### 4. Jalankan Services dengan Docker Compose

```bash
docker-compose up -d
```

Services yang akan berjalan:
- MySQL
- Elasticsearch
- Apache Kafka

Tunggu beberapa detik sampai services fully healthy.

### 5. Run Database Migrations

```bash
# Apply all migrations (up)
migrate -path migrations -database "mysql://root:root@tcp(localhost:3306)/gocommerce" up

# Rollback semua migrations (down)
migrate -path migrations -database "mysql://root:root@tcp(localhost:3306)/gocommerce" down
```

## ⚙️ Konfigurasi

### Database Configuration

Aplikasi menggunakan MySQL dengan connection pool yang dapat dikonfigurasi:

```env
DATABASE_HOST=localhost        # Host MySQL
DATABASE_PORT=3306           # Port MySQL
DATABASE_NAME=gocommerce     # Database name
DATABASE_USER=root           # Username
DATABASE_PASSWORD=root       # Password

POOL_MAX_IDLE_CONNS=10       # Max idle connections
POOL_MAX_OPEN_CONNS=100      # Max open connections
POOL_CONN_IDLE_TIME=600      # Connection idle timeout (seconds)
POOL_CONN_LIFE_TIME=3600     # Connection lifetime (seconds)
```

### Elasticsearch Configuration

Untuk pencarian full-text:

```env
ELASTICSEARCH_ADDRESSES=http://localhost:9201
ELASTICSEARCH_USER_INDEX=users
ELASTICSEARCH_PRODUCT_INDEX=products
ELASTICSEARCH_ORDER_INDEX=orders
```

Elasticsearch secara otomatis membuat indices dengan mappings yang sudah ditentukan di folder `internal/infrastructure/elasticsearch/mappings/`.

### Kafka Configuration

Untuk event streaming dan asynchronous processing:

```env
KAFKA_GROUP_ID=gocommerce-group
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=order-events
```

## 🎯 Menjalankan Aplikasi

```bash
# Terminal 1: Jalankan API Server
cd cmd/api
go run main.go

# Terminal 2 (optional): Jalankan Worker/Consumer
cd cmd/worker
go run main.go
```

API akan tersedia di `http://localhost:3000`

## 📁 Struktur Project

```
gocommerce/
├── cmd/                           # Entry points
│   ├── api/                       # API Server
│   │   └── main.go
│   └── worker/                    # Background Worker/Consumer
│       └── main.go
│
├── config/                        # Configuration
│   └── config.go                 # Load & validate environment variables
│
├── internal/                      # Internal packages
│   ├── bootstrap.go              # Initialize routes & modules
│   ├── infrastructure/           # Infrastructure layer
│   │   ├── database.go           # Database connections
│   │   ├── fiber.go              # Fiber app setup
│   │   ├── logger.go             # Logging setup
│   │   ├── elasticsearch/        # Elasticsearch client & integration
│   │   │   ├── client.go
│   │   │   ├── indices.go
│   │   │   └── mappings/         # Index mappings
│   │   │       ├── users_mapping.json
│   │   │       ├── products_mapping.json
│   │   │       └── orders_mapping.json
│   │   └── kafka/                # Kafka producer & consumer
│   │       ├── producer.go
│   │       └── consumer.go
│   │
│   └── modules/                  # Domain modules (DDD)
│       ├── user/                 # User module
│       │   ├── module.go         # Module registration
│       │   ├── delivery/         # HTTP handlers
│       │   ├── dto/              # Data Transfer Objects
│       │   ├── entity/           # Domain entities
│       │   ├── repository/       # Data access layer
│       │   ├── search/           # Elasticsearch operations
│       │   └── usecase/          # Business logic
│       │
│       ├── product/              # Product module
│       │   ├── module.go
│       │   ├── delivery/
│       │   ├── dto/
│       │   ├── entity/
│       │   ├── repository/
│       │   ├── search/           # Product search with Elasticsearch
│       │   └── usecase/
│       │
│       ├── order/                # Order module
│       │   ├── module.go
│       │   ├── delivery/
│       │   │   ├── http/         # HTTP handlers
│       │   │   └── messaging/    # Kafka event handlers
│       │   ├── dto/
│       │   ├── entity/
│       │   ├── repository/
│       │   ├── search/           # Order search
│       │   └── usecase/
│       │
│       └── shared/               # Shared utilities
│           ├── pagination_filter.go
│           ├── repository.go     # Base repository interface
│           └── file_upload/      # File upload service
│               ├── config.go
│               ├── service.go
│               └── validator.go
│
├── migrations/                    # Database migrations
│   ├── 20260511042410_create_users_table.up.sql
│   ├── 20260511042410_create_users_table.down.sql
│   ├── 20260511043106_create_products_table.up.sql
│   ├── 20260511043106_create_products_table.down.sql
│   ├── 20260511043633_create_orders_table.up.sql
│   ├── 20260511043633_create_orders_table.down.sql
│   ├── 20260511044442_create_order_items_table.up.sql
│   └── 20260511044442_create_order_items_table.down.sql
│
├── pkg/                          # Public packages
│   ├── response/                 # Standard response format
│   ├── validation/               # Request validation
│   └── storage/                  # Storage utilities
│
├── uploads/                      # Upload directory
│   └── products/                 # Product images
│
├── docker-compose.yml            # Docker services configuration
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
└── README.md                     # This file
```

## 📚 Architecture Overview

### Layered Architecture

```
┌─────────────────────────────────┐
│        HTTP Handlers            │  (Delivery Layer)
│        (REST API)               │
└─────────────────────────────────┘
              ↓
┌─────────────────────────────────┐
│      Use Cases / Services       │  (Business Logic Layer)
│      (Business Logic)           │
└─────────────────────────────────┘
              ↓
┌─────────────────────────────────┐
│      Repository Pattern         │  (Data Access Layer)
│   (Database & Search Queries)   │
└─────────────────────────────────┘
              ↓
┌─────────────────────────────────┐
│    Infrastructure Layer         │  (External Services)
│  (DB, Elasticsearch, Kafka)     │
└─────────────────────────────────┘
```

### Module Structure

Setiap module mengikuti Domain-Driven Design dengan struktur:

- **Entity**: Domain objects
- **DTO**: Data Transfer Objects untuk request/response
- **Repository**: Data access abstraction
- **UseCase**: Business logic orchestration
- **Delivery**: Handler untuk HTTP requests
- **Search**: Elasticsearch integration

## 🔌 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### User Endpoints
```
GET    /users              # Get all users
GET    /users/:id          # Get user by ID
POST   /users              # Create new user
PUT    /users/:id          # Update user
DELETE /users/:id          # Delete user
```

### Product Endpoints
```
GET    /products           # Get all products (dengan pagination)
GET    /products/:id       # Get product by ID
POST   /products           # Create new product
PUT    /products/:id       # Update product
DELETE /products/:id       # Delete product
POST   /products/search    # Search products (full-text search)
PUT    /products/:id/stock # Update product stock
POST   /products/images    # Upload product image
```

### Order Endpoints
```
GET    /orders             # Get all orders
GET    /orders/:id         # Get order by ID
POST   /orders             # Create new order
PUT    /orders/:id         # Update order status
GET    /orders/search      # Search orders
```

## 🔍 File Upload Configuration

File upload dikelola di folder `internal/shared/file_upload/` dengan support untuk:

- Validasi tipe file (image types)
- Validasi ukuran file
- Penyimpanan di folder `uploads/products/`
- UUID-based filename untuk keamanan

## 🔍 Elasticsearch Integration

### Indices

Aplikasi membuat 3 indices Elasticsearch:

1. **users**: Index untuk user data
2. **products**: Index untuk product data dengan full-text search
3. **orders**: Index untuk order data

Setiap index memiliki mapping yang sudah dikonfigurasi untuk optimal search experience.

### Searching

```go
// Contoh: Search products
GET /api/v1/products/search?q=laptop&page=1&limit=10

// Response akan mencari di Elasticsearch dan return hasil yang relevan
```

## 📡 Kafka Integration

### Event Streaming

Aplikasi menggunakan Kafka untuk event-driven architecture:

- **Topic**: `order-events`
- **Group ID**: `gocommerce-group`
- **Events**: Order creation, status updates, payment notifications

### Producer (API Server)
- Mengirim events saat order dibuat atau diupdate

### Consumer (Worker)
- Mendengarkan events dan melakukan processing asynchronously
- Update database, send notifications, etc.

Jalankan worker dengan:
```bash
go run cmd/worker/main.go
```

## 👤 Author

**Alfian Yulianto**
- GitHub: [@alfianyulianto](https://github.com/alfianyulianto)

