# 🔒 Security Documentation

## 🛡️ Security Overview

Aplikasi Online Exam mengimplementasikan multiple layer security untuk melindungi data dan pengguna.

## 🔐 Authentication & Authorization

### JWT (JSON Web Token)
```javascript
// Token Structure
{
  "user_id": 123,
  "email": "user@example.com",
  "role": "user",
  "exp": 1642675200, // 24 hours
  "iat": 1642588800
}
```

**Security Features:**
- Token expiration: 24 hours
- Secure secret key (environment variable)
- Role-based access control
- Automatic token validation

### Password Security
```go
// Password Hashing with bcrypt
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

**Security Features:**
- bcrypt hashing with salt
- Cost factor: 12 (default)
- One-way encryption
- Salt generation per password

## 🌐 Network Security

### HTTPS/TLS
```nginx
# Nginx SSL Configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
ssl_prefer_server_ciphers off;
```

**Security Features:**
- TLS 1.2+ only
- Strong cipher suites
- HSTS headers
- SSL certificate validation

### CORS Configuration
```go
// CORS Middleware
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))
```

## 🗄️ Database Security

### SQL Injection Prevention
```go
// Parameterized Queries
db.Where("email = ?", email).First(&user)
```

**Security Features:**
- Parameterized queries only
- Input validation
- Type checking
- ORM protection (GORM)

### Database Access Control
```sql
-- User Permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON answers TO app_user;
GRANT SELECT ON questions TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON users TO admin_user;
```

## 🔍 Input Validation

### Server-side Validation
```go
// Email Validation
func validateEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}

// Password Validation
func validatePassword(password string) bool {
    return len(password) >= 8
}
```

### Frontend Validation
```javascript
// Client-side Validation
const validateForm = (data) => {
    const errors = {};
    
    if (!data.email || !emailRegex.test(data.email)) {
        errors.email = 'Email tidak valid';
    }
    
    if (!data.password || data.password.length < 8) {
        errors.password = 'Password minimal 8 karakter';
    }
    
    return errors;
};
```

## 🚫 Rate Limiting

### API Rate Limiting
```go
// Rate Limiting Middleware
app.Use(limiter.New(limiter.Config{
    Max:        100,  // requests
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
}))
```

**Limits:**
- Login attempts: 5 per minute
- API requests: 100 per minute
- Exam submissions: 10 per hour

## 📊 Session Management

### Redis Session Storage
```go
// Session Configuration
store := session.New(session.Config{
    Storage:    redisClient,
    Expiration: 24 * time.Hour,
    KeyLookup:  "cookie:session_id",
})
```

**Security Features:**
- Secure session storage
- Automatic expiration
- Session invalidation on logout
- CSRF protection

## 🔐 Access Control

### Role-based Access Control (RBAC)
```go
// Admin Middleware
func adminMiddleware(c *fiber.Ctx) error {
    role := c.Locals("role")
    if role != "admin" {
        return c.Status(403).JSON(fiber.Map{
            "success": false,
            "message": "Akses admin diperlukan",
        })
    }
    return c.Next()
}
```

**Roles:**
- `user`: Peserta ujian
- `admin`: Administrator sistem

### Route Protection
```javascript
// Frontend Route Protection
const ProtectedRoute = ({ children, requiredRole }) => {
    const { user } = useAuth();
    
    if (!user) {
        return <Navigate to="/login" />;
    }
    
    if (requiredRole && user.role !== requiredRole) {
        return <Navigate to="/dashboard" />;
    }
    
    return children;
};
```

## 🛡️ Data Protection

### Sensitive Data Handling
```go
// Password Masking in Logs
func maskPassword(password string) string {
    if len(password) <= 2 {
        return "***"
    }
    return password[:1] + "***" + password[len(password)-1:]
}
```

### Data Encryption
- Passwords: bcrypt hashing
- Tokens: JWT signing
- Database: TLS connection
- File uploads: Secure storage

## 🔍 Security Monitoring

### Logging
```go
// Security Event Logging
func logSecurityEvent(event string, userID uint, details map[string]interface{}) {
    log.Printf("SECURITY: %s | User: %d | Details: %v", event, userID, details)
}
```

**Logged Events:**
- Login attempts (success/failure)
- Admin actions
- Data modifications
- Security violations

### Error Handling
```go
// Secure Error Responses
func handleError(c *fiber.Ctx, err error) error {
    // Don't expose internal errors
    log.Printf("Error: %v", err)
    
    return c.Status(500).JSON(fiber.Map{
        "success": false,
        "message": "Terjadi kesalahan sistem",
    })
}
```

## 🚨 Incident Response

### Security Breach Response
1. **Immediate Actions:**
   - Isolate affected systems
   - Preserve evidence
   - Notify stakeholders

2. **Investigation:**
   - Analyze logs
   - Identify root cause
   - Assess impact

3. **Recovery:**
   - Patch vulnerabilities
   - Restore from backups
   - Update security measures

### Contact Information
- Security Team: security@company.com
- Emergency: +62-xxx-xxx-xxxx
- Bug Reports: GitHub Issues

## 📋 Security Checklist

### Development
- [ ] Input validation implemented
- [ ] SQL injection prevention
- [ ] XSS protection
- [ ] CSRF protection
- [ ] Secure headers configured

### Deployment
- [ ] HTTPS enabled
- [ ] SSL certificate valid
- [ ] Firewall configured
- [ ] Database access restricted
- [ ] Backup encryption enabled

### Monitoring
- [ ] Security logs enabled
- [ ] Error tracking configured
- [ ] Performance monitoring active
- [ ] Alert system operational

## 🔄 Security Updates

### Regular Updates
- Dependencies: Monthly
- Security patches: As needed
- SSL certificates: 90 days
- Password policies: Quarterly

### Vulnerability Management
- Automated scanning
- Manual code review
- Third-party audits
- Penetration testing

## 📚 Security Resources

### Documentation
- OWASP Top 10
- Go Security Best Practices
- React Security Guidelines
- PostgreSQL Security

### Tools
- OWASP ZAP (penetration testing)
- SonarQube (code analysis)
- Snyk (dependency scanning)
- Let's Encrypt (SSL certificates) 