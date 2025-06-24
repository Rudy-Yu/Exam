# 🗄️ Database Documentation

## 📊 Database Schema Overview

### Entity Relationship Diagram
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Users    │    │  Questions  │    │   Answers   │
│             │    │             │    │             │
│ id (PK)     │    │ id (PK)     │    │ id (PK)     │
│ email       │    │ exam_id     │    │ participant_id (FK)
│ password    │    │ question_text│   │ question_id (FK)
│ role        │    │ correct_answer│  │ answer_text │
│ created_at  │    │ weight      │    │ submitted_at│
│ updated_at  │    │ created_at  │    │ is_draft    │
└─────────────┘    └─────────────┘    └─────────────┘
```

## 📋 Table Definitions

### 1. Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Columns:**
- `id`: Primary key, auto-increment
- `email`: Unique email address
- `password`: Hashed password (bcrypt)
- `role`: User role ('user' or 'admin')
- `created_at`: Record creation timestamp
- `updated_at`: Record update timestamp

**Indexes:**
- Primary key on `id`
- Unique index on `email`
- Index on `role` for role-based queries

### 2. Questions Table
```sql
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    exam_id INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    weight INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Columns:**
- `id`: Primary key, auto-increment
- `exam_id`: Foreign key to exams table
- `question_text`: Question content
- `correct_answer`: Correct answer
- `weight`: Question weight for scoring
- `created_at`: Record creation timestamp
- `updated_at`: Record update timestamp

**Indexes:**
- Primary key on `id`
- Index on `exam_id` for exam-based queries
- Index on `weight` for scoring calculations

### 3. Answers Table
```sql
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    participant_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    answer_text TEXT NOT NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_draft BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (participant_id) REFERENCES users(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);
```

**Columns:**
- `id`: Primary key, auto-increment
- `participant_id`: Foreign key to users table
- `question_id`: Foreign key to questions table
- `answer_text`: User's answer
- `submitted_at`: Answer submission timestamp
- `is_draft`: Flag for draft answers

**Indexes:**
- Primary key on `id`
- Foreign key indexes on `participant_id` and `question_id`
- Composite index on `(participant_id, question_id)`
- Index on `is_draft` for draft filtering

## 🔗 Relationships

### One-to-Many Relationships
1. **User → Answers**: One user can have multiple answers
2. **Question → Answers**: One question can have multiple answers from different users

### Many-to-One Relationships
1. **Answers → User**: Multiple answers belong to one user
2. **Answers → Question**: Multiple answers belong to one question

## 📊 Data Types & Constraints

### String Fields
- `email`: VARCHAR(255) - Email addresses
- `password`: VARCHAR(255) - Hashed passwords
- `role`: VARCHAR(50) - User roles
- `question_text`: TEXT - Question content
- `correct_answer`: VARCHAR(255) - Correct answers
- `answer_text`: TEXT - User answers

### Numeric Fields
- `id`: SERIAL - Auto-incrementing primary keys
- `exam_id`: INTEGER - Exam identifiers
- `weight`: INTEGER - Question weights
- `participant_id`: INTEGER - User identifiers
- `question_id`: INTEGER - Question identifiers

### Boolean Fields
- `is_draft`: BOOLEAN - Draft answer flag

### Timestamp Fields
- `created_at`: TIMESTAMP - Record creation time
- `updated_at`: TIMESTAMP - Record update time
- `submitted_at`: TIMESTAMP - Answer submission time

## 🔍 Query Examples

### 1. Get User with Answers
```sql
SELECT u.email, q.question_text, a.answer_text, a.submitted_at
FROM users u
JOIN answers a ON u.id = a.participant_id
JOIN questions q ON a.question_id = q.id
WHERE u.id = $1 AND a.is_draft = FALSE
ORDER BY a.submitted_at DESC;
```

### 2. Get Exam Results
```sql
SELECT 
    u.email,
    COUNT(CASE WHEN a.answer_text = q.correct_answer THEN 1 END) as correct_answers,
    COUNT(*) as total_questions,
    SUM(CASE WHEN a.answer_text = q.correct_answer THEN q.weight ELSE 0 END) as score
FROM users u
JOIN answers a ON u.id = a.participant_id
JOIN questions q ON a.question_id = q.id
WHERE q.exam_id = $1 AND a.is_draft = FALSE
GROUP BY u.id, u.email;
```

### 3. Get Draft Answers
```sql
SELECT q.question_text, a.answer_text
FROM answers a
JOIN questions q ON a.question_id = q.id
WHERE a.participant_id = $1 AND a.is_draft = TRUE;
```

## 🔧 Database Configuration

### Connection Pool Settings
```go
// PostgreSQL connection pool
MaxConns: 20
MinConns: 5
MaxConnLifetime: 1 hour
MaxConnIdleTime: 30 minutes
```

### Redis Configuration
```go
// Redis connection pool
PoolSize: 100
MinIdleConns: 10
MaxRetries: 3
```

## 📈 Performance Optimization

### Indexing Strategy
1. **Primary Keys**: All tables have auto-increment primary keys
2. **Foreign Keys**: Indexed for join performance
3. **Composite Indexes**: For common query patterns
4. **Partial Indexes**: For filtered queries

### Query Optimization
1. **Selective Columns**: Only fetch required columns
2. **Limit Results**: Use LIMIT for large result sets
3. **Efficient Joins**: Use appropriate join types
4. **Caching**: Redis for frequently accessed data

## 🔒 Security Considerations

### Data Protection
1. **Password Hashing**: bcrypt with salt
2. **Input Validation**: Server-side validation
3. **SQL Injection Prevention**: Parameterized queries
4. **Access Control**: Role-based permissions

### Backup Strategy
1. **Daily Backups**: Automated database backups
2. **Point-in-Time Recovery**: WAL archiving
3. **Backup Verification**: Regular restore tests
4. **Offsite Storage**: Secure backup storage

## 🛠️ Maintenance

### Regular Tasks
1. **VACUUM**: Clean up dead tuples
2. **ANALYZE**: Update statistics
3. **REINDEX**: Rebuild indexes
4. **Log Rotation**: Manage log files

### Monitoring
1. **Connection Count**: Monitor active connections
2. **Query Performance**: Slow query analysis
3. **Disk Usage**: Monitor table sizes
4. **Index Usage**: Track index effectiveness 