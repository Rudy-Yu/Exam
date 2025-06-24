# 🚀 Deployment Documentation

## 📋 Deployment Overview

Panduan lengkap untuk deployment aplikasi Online Exam ke production environment.

## 🏗️ Infrastructure Requirements

### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 50GB SSD
- **OS**: Ubuntu 20.04 LTS
- **Network**: Public IP with domain

### Recommended Requirements
- **CPU**: 4 cores
- **RAM**: 8GB
- **Storage**: 100GB SSD
- **OS**: Ubuntu 22.04 LTS
- **Network**: Load balancer + CDN

## 🔧 Server Setup

### 1. Initial Server Configuration
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y curl wget git unzip software-properties-common

# Create deployment user
sudo adduser deploy
sudo usermod -aG sudo deploy
```

### 2. Install Dependencies

#### Install Go
```bash
# Download Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### Install Node.js
```bash
# Install Node.js 18
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

#### Install PostgreSQL
```bash
# Install PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Start and enable service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql
CREATE DATABASE examdb;
CREATE USER examuser WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE examdb TO examuser;
\q
```

#### Install Redis
```bash
# Install Redis
sudo apt install -y redis-server

# Configure Redis
sudo nano /etc/redis/redis.conf
# Set: requirepass your_redis_password

# Start and enable service
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

#### Install Nginx
```bash
# Install Nginx
sudo apt install -y nginx

# Start and enable service
sudo systemctl start nginx
sudo systemctl enable nginx
```

## 🔐 SSL Certificate Setup

### Using Let's Encrypt
```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 📁 Application Deployment

### 1. Clone Repository
```bash
# Switch to deploy user
sudo su - deploy

# Clone repository
git clone https://github.com/username/online-exam-app.git
cd online-exam-app
```

### 2. Backend Deployment
```bash
# Navigate to backend
cd backend

# Install dependencies
go mod download

# Build application
go build -o online-exam-app

# Create systemd service
sudo nano /etc/systemd/system/online-exam.service
```

**Service Configuration:**
```ini
[Unit]
Description=Online Exam Backend
After=network.target postgresql.service redis-server.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/online-exam-app/backend
Environment=GO_ENV=production
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=examuser
Environment=DB_PASS=secure_password
Environment=DB_NAME=examdb
Environment=REDIS_HOST=localhost
Environment=REDIS_PORT=6379
Environment=REDIS_PASS=your_redis_password
Environment=JWT_SECRET=your_jwt_secret_key
ExecStart=/home/deploy/online-exam-app/backend/online-exam-app
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### 3. Frontend Deployment
```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Copy to web directory
sudo cp -r build/* /var/www/html/
sudo chown -R www-data:www-data /var/www/html/
```

## 🌐 Nginx Configuration

### Main Configuration
```nginx
# /etc/nginx/sites-available/online-exam
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    # Frontend
    location / {
        root /var/www/html;
        try_files $uri $uri/ /index.html;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # API Proxy
    location /api {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Rate limiting
        limit_req zone=api burst=20 nodelay;
    }

    # Health check
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

### Rate Limiting Configuration
```nginx
# /etc/nginx/nginx.conf
http {
    # Rate limiting zones
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=1r/s;
    
    # Include site configurations
    include /etc/nginx/sites-enabled/*;
}
```

## 🔄 Deployment Process

### 1. Start Services
```bash
# Enable and start backend service
sudo systemctl enable online-exam
sudo systemctl start online-exam

# Enable and start Nginx
sudo systemctl enable nginx
sudo systemctl start nginx

# Check status
sudo systemctl status online-exam
sudo systemctl status nginx
```

### 2. Database Migration
```bash
# Run database migrations
cd /home/deploy/online-exam-app/backend
go run main.go migrate
```

### 3. Verify Deployment
```bash
# Check application logs
sudo journalctl -u online-exam -f

# Check Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# Test endpoints
curl -k https://yourdomain.com/health
curl -k https://yourdomain.com/api/exams
```

## 📊 Monitoring Setup

### 1. Application Monitoring
```bash
# Install monitoring tools
sudo apt install -y htop iotop nethogs

# Setup log rotation
sudo nano /etc/logrotate.d/online-exam
```

### 2. Database Monitoring
```bash
# Enable PostgreSQL logging
sudo nano /etc/postgresql/14/main/postgresql.conf
# Set: log_statement = 'all'
# Set: log_destination = 'stderr'
# Set: logging_collector = on

# Restart PostgreSQL
sudo systemctl restart postgresql
```

### 3. Redis Monitoring
```bash
# Monitor Redis
redis-cli monitor
redis-cli info memory
```

## 🔧 Maintenance Procedures

### 1. Application Updates
```bash
# Pull latest code
cd /home/deploy/online-exam-app
git pull origin main

# Rebuild backend
cd backend
go build -o online-exam-app

# Restart service
sudo systemctl restart online-exam

# Rebuild frontend
cd ../frontend
npm install
npm run build
sudo cp -r build/* /var/www/html/
```

### 2. Database Backup
```bash
# Create backup script
sudo nano /usr/local/bin/backup-db.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="examdb"

mkdir -p $BACKUP_DIR
pg_dump $DB_NAME > $BACKUP_DIR/${DB_NAME}_${DATE}.sql

# Keep only last 7 days
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
```

### 3. SSL Certificate Renewal
```bash
# Test renewal
sudo certbot renew --dry-run

# Manual renewal
sudo certbot renew
sudo systemctl reload nginx
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Application Won't Start
```bash
# Check logs
sudo journalctl -u online-exam -n 50

# Check environment variables
sudo systemctl show online-exam --property=Environment
```

#### 2. Database Connection Issues
```bash
# Test database connection
psql -h localhost -U examuser -d examdb

# Check PostgreSQL status
sudo systemctl status postgresql
```

#### 3. Redis Connection Issues
```bash
# Test Redis connection
redis-cli ping

# Check Redis status
sudo systemctl status redis-server
```

#### 4. Nginx Issues
```bash
# Test configuration
sudo nginx -t

# Check Nginx status
sudo systemctl status nginx
```

## 📈 Performance Optimization

### 1. Database Optimization
```sql
-- Analyze tables
ANALYZE users;
ANALYZE questions;
ANALYZE answers;

-- Create indexes
CREATE INDEX idx_answers_participant_question ON answers(participant_id, question_id);
CREATE INDEX idx_questions_exam_id ON questions(exam_id);
```

### 2. Application Optimization
```bash
# Enable Go profiling
export GODEBUG=gctrace=1

# Monitor memory usage
go tool pprof http://localhost:6060/debug/pprof/heap
```

### 3. Nginx Optimization
```nginx
# Enable gzip compression
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# Enable caching
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 🔒 Security Hardening

### 1. Firewall Configuration
```bash
# Install UFW
sudo apt install -y ufw

# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### 2. Fail2ban Setup
```bash
# Install Fail2ban
sudo apt install -y fail2ban

# Configure for Nginx
sudo nano /etc/fail2ban/jail.local
```

### 3. Regular Security Updates
```bash
# Setup automatic updates
sudo apt install -y unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

## 📞 Support Information

### Contact Details
- **System Administrator**: admin@company.com
- **Emergency Contact**: +62-xxx-xxx-xxxx
- **Documentation**: https://docs.company.com

### Monitoring Dashboard
- **Application**: https://yourdomain.com/health
- **Server Status**: https://status.company.com
- **Logs**: /var/log/online-exam/ 