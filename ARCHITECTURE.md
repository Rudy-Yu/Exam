# 🏗️ Architecture Documentation

## 📐 System Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (React.js)    │◄──►│   (Go/Fiber)    │◄──►│  (PostgreSQL)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │     Redis       │
                       │   (Cache/Session)│
                       └─────────────────┘
```

## 🔧 Technology Stack

### Frontend Layer
- **Framework**: React.js 18+
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **State Management**: React Hooks
- **Build Tool**: Create React App

### Backend Layer
- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **ORM**: GORM v2
- **Authentication**: JWT
- **Session**: Redis

### Data Layer
- **Primary DB**: PostgreSQL 14+
- **Cache**: Redis 6+
- **Connection Pooling**: pgxpool

### Infrastructure
- **Reverse Proxy**: Nginx (Production)
- **SSL/TLS**: Let's Encrypt
- **Process Manager**: systemd

## 🏛️ Architecture Patterns

### 1. **Layered Architecture**
```
┌─────────────────────────────────────┐
│           Presentation Layer        │ ← React Components
├─────────────────────────────────────┤
│           Business Logic Layer      │ ← Go Handlers
├─────────────────────────────────────┤
│           Data Access Layer         │ ← GORM Models
├─────────────────────────────────────┤
│           Data Storage Layer        │ ← PostgreSQL + Redis
└─────────────────────────────────────┘
```

### 2. **RESTful API Design**
- Resource-based URLs
- HTTP methods (GET, POST, PUT, DELETE)
- JSON request/response
- Stateless communication

### 3. **Security Patterns**
- JWT-based authentication
- Role-based access control (RBAC)
- CORS configuration
- Input validation
- SQL injection prevention

## 📊 Data Flow

### 1. **User Authentication Flow**
```
User Login → JWT Token → Redis Session → Protected Routes
```

### 2. **Exam Flow**
```
Start Exam → Redis Session → Timer Sync → Auto-save → Submit → PostgreSQL
```

### 3. **Admin Flow**
```
Admin Login → JWT Token → CRUD Operations → Database → Response
```

## 🔒 Security Architecture

### Authentication
- JWT tokens with 24-hour expiration
- Secure token storage in localStorage
- Automatic token refresh mechanism

### Authorization
- Role-based access (user/admin)
- Route protection middleware
- API endpoint security

### Data Protection
- Password hashing (bcrypt)
- HTTPS/TLS encryption
- Input sanitization
- SQL injection prevention

## 📈 Performance Architecture

### Caching Strategy
- Redis for session storage
- Redis for auto-save drafts
- Connection pooling for database

### Scalability
- Horizontal scaling ready
- Stateless backend design
- Load balancer compatible

### Monitoring
- Application logs
- Database performance metrics
- Redis monitoring
- Error tracking

## 🚀 Deployment Architecture

### Development Environment
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   React     │    │    Go       │    │ PostgreSQL  │
│  :3001      │    │   :3000     │    │   :5433     │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Production Environment
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Nginx     │    │    Go       │    │ PostgreSQL  │
│  :80/443    │    │   :3000     │    │   :5432     │
└─────────────┘    └─────────────┘    └─────────────┘
                              │
                              ▼
                       ┌─────────────┐
                       │    Redis    │
                       │   :6379     │
                       └─────────────┘
```

## 🔄 System Integration

### External Dependencies
- PostgreSQL Database
- Redis Cache
- SSL Certificate Provider (Let's Encrypt)

### Internal Dependencies
- Go modules
- npm packages
- System libraries

## 📋 Configuration Management

### Environment Variables
- Database configuration
- Redis configuration
- JWT secrets
- SSL certificates
- API endpoints

### Configuration Files
- `go.mod` (Go dependencies)
- `package.json` (Node.js dependencies)
- `nginx.conf` (Web server config)
- `.env` files (Environment variables)

## 🔍 Monitoring & Logging

### Application Logs
- Request/response logging
- Error tracking
- Performance metrics
- Security events

### Infrastructure Monitoring
- Server health checks
- Database performance
- Redis memory usage
- Network connectivity

## 🛡️ Disaster Recovery

### Backup Strategy
- Database backups (daily)
- Configuration backups
- Code repository backups

### Recovery Procedures
- Database restoration
- Application restart
- SSL certificate renewal
- Configuration recovery 