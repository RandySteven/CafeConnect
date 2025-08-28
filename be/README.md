# CafeConnect Backend

A robust backend service for the CafeConnect mobile application that helps users discover nearby cafes based on theme, atmosphere, and price. Built with Go using clean architecture principles.

![CafeConnect Logo](cafeConnect-logo.png)

## 🏗️ Architecture Overview

CafeConnect Backend follows a clean architecture pattern with the following layers:

```
┌─────────────────────────────────────────────────────────────┐
│                        Handlers (APIs)                     │
├─────────────────────────────────────────────────────────────┤
│                       Use Cases                            │
├─────────────────────────────────────────────────────────────┤
│                     Repositories                           │
├─────────────────────────────────────────────────────────────┤
│                        Database                            │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

- **Handlers**: HTTP API endpoints with request/response handling
- **Use Cases**: Business logic and orchestration
- **Repositories**: Data access layer with database operations
- **Entities**: Domain models and data structures
- **Interfaces**: Contract definitions for dependency injection
- **Middlewares**: Authentication, logging, CORS, rate limiting
- **Pkg**: External service integrations (Midtrans, Kafka, etc.)

## 🚀 Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gorilla Mux (HTTP routing)
- **Database**: MySQL 8.0
- **Cache**: Redis 8.0
- **Message Broker**: Apache Kafka 3.9.1
- **Search Engine**: Elasticsearch 9.0.2
- **Payment Gateway**: Midtrans
- **Cloud Storage**: Google Cloud Storage
- **Authentication**: JWT + OAuth2 (Google)

### Development Tools
- **Containerization**: Docker & Docker Compose
- **Build Tool**: Make
- **Testing**: Go testing framework
- **Mocking**: Mockery

## 📁 Project Structure

```
├── cmd/                    # Application entry points
│   └── main/
│       ├── http/          # HTTP server
│       ├── migration/     # Database migrations
│       ├── seed/          # Database seeding
│       └── consumer/      # Kafka consumers
├── apps/                  # Application configuration
├── configs/               # Configuration management
├── entities/              # Domain models and DTOs
├── handlers/              # HTTP handlers and middleware
├── interfaces/            # Interface definitions
├── middlewares/           # HTTP middleware
├── pkg/                   # External service packages
├── repositories/          # Data access layer
├── routes/                # Route definitions
├── usecases/              # Business logic
├── utils/                 # Utility functions
└── caches/                # Cache implementations
```

## 🔌 API Endpoints

### Base URL
```
http://localhost:8081
```

### 1. Development & Health
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/dev/check-health` | Health check endpoint | ❌ |

### 2. Authentication & Onboarding
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/onboarding/register` | User registration | ❌ |
| `POST` | `/onboarding/login` | User login | ❌ |
| `GET` | `/onboarding` | Get user profile | ✅ |
| `GET` | `/onboarding/verify-token/{token}` | Verify email token | ❌ |

### 3. User Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/roles` | Get available roles | ❌ |

### 4. Address Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/addresses` | Add user address | ✅ |
| `GET` | `/addresses` | Get user addresses | ✅ |

### 5. Cafe Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/cafes/franchises` | Get cafe franchises | ❌ |
| `POST` | `/cafes/franchises` | Register cafe & franchise | ❌ |
| `GET` | `/cafes/{id}` | Get cafe details | ❌ |
| `POST` | `/cafes` | Get cafes by radius | ❌ |
| `POST` | `/cafes/add-outlet` | Add cafe outlet | ❌ |

### 6. Product Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/products` | Add new product | ❌ |
| `GET` | `/products/{id}` | Get product details | ❌ |
| `POST` | `/products/cafe-list` | Get products by cafe | ❌ |

### 7. Review System
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/reviews/cafe-review` | Get cafe reviews | ❌ |
| `POST` | `/reviews` | Add cafe review | ✅ |

### 8. Shopping Cart
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/carts` | Add item to cart | ✅ |
| `GET` | `/carts` | Get user cart | ✅ |

### 9. Transactions
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/transactions/v1/check-out` | Checkout V1 | ✅ |
| `POST` | `/transactions/v2/check-out` | Checkout V2 | ✅ |
| `GET` | `/transactions/{transactionCode}` | Get transaction by code | ✅ |
| `GET` | `/transactions` | Get user transactions | ✅ |
| `POST` | `/transactions/receipt` | Check transaction receipt | ✅ |
| `POST` | `/transactions/payment` | Payment confirmation | ✅ |

## 🔐 Authentication

The API uses JWT tokens for authentication. Protected endpoints require the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

## 🗄️ Database Schema

### Core Tables
- **users**: User accounts and profiles
- **cafes**: Cafe information and locations
- **cafe_franchises**: Franchise management
- **products**: Product catalog
- **cafe_products**: Cafe-specific product offerings
- **transactions**: Order transactions
- **transaction_details**: Order line items
- **carts**: Shopping cart items
- **reviews**: User reviews and ratings
- **addresses**: User delivery addresses

## 🔄 Data Flow

### 1. User Registration Flow
```
User Input → Validation → Password Hash → Database Save → Email Verification
```

### 2. Cafe Discovery Flow
```
Location Input → Radius Search → Filter by Criteria → Return Results
```

### 3. Order Processing Flow
```
Cart Items → Transaction Creation → Payment Processing → Order Confirmation
```

### 4. Payment Flow
```
Order → Midtrans Integration → Payment Gateway → Webhook → Status Update
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional, for convenience)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd CafeConnect/be
   ```

2. **Set up environment variables**
   ```bash
   cp files/env/.env.example files/env/.env
   # Edit .env with your configuration
   ```

3. **Run with Docker (Recommended)**
   ```bash
   # Start all services
   docker compose up --build
   
   # Run in background
   docker compose up --build -d
   
   # Refresh containers
   docker compose up --build -d
   ```

4. **Run locally**
   ```bash
   # Install dependencies
   go mod download
   
   # Run the application
   make run
   
   # Or directly
   go run cmd/main/http/main.go -config files/yml/cofeConnect.local.yaml
   ```

### Database Operations

```bash
# Run migrations
make migration

# Seed database
make seed

# Drop database
make drop

# Refresh database (drop + migrate + seed)
make refresh
```

### Additional Commands

```bash
# Run tests with coverage
make test

# Generate mocks
make mockery

# Build binary
make build

# Run consumer
make consumer

# Run scheduler
make scheduler
```

## 🔧 Configuration

The application uses YAML configuration files located in `files/yml/`. Key configuration sections:

- **Server**: Host, port, timeouts
- **Database**: MySQL connection settings
- **Redis**: Cache configuration
- **Kafka**: Message broker settings
- **Midtrans**: Payment gateway configuration
- **Google Cloud**: Storage and OAuth settings

## 📊 Monitoring & Logging

- **Health Checks**: `/dev/check-health` endpoint
- **Request Logging**: Automatic request/response logging
- **Error Tracking**: Structured error logging with request IDs
- **Performance**: Request timeout and rate limiting

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./handlers/...

# Run with verbose output
go test -v ./...
```

## 📦 Deployment

### Docker Deployment
```bash
# Build and run
docker compose up --build

# Production build
docker build -t cafeconnect-backend .
```

### Environment Variables
- `ENV`: Environment (development/staging/production)
- `DB_HOST`: Database host
- `REDIS_HOST`: Redis host
- `KAFKA_BROKERS`: Kafka broker addresses
- `MIDTRANS_SERVER_KEY`: Midtrans server key
- `GOOGLE_CLOUD_PROJECT_ID`: Google Cloud project ID

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📝 License

This project is proprietary software. All rights reserved.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation and examples

---

**CafeConnect Backend** - Powering the future of cafe discovery and ordering.
