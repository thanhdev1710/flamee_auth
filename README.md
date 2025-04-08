# 🔥 Flamee Auth Service

Authentication service for the **Flamee** social network project.  
Built with **Golang**, using **Gin**, **JWT**, and **PostgreSQL**.

---

## 📌 Features

- User registration & login
- JWT access/refresh token
- Email verification
- Password reset
- Rate limiting per route & IP
- Custom logger using zap + lumberjack
- Secured with CORS, Helmet headers, API key middleware
- Modular structure with service & repository layers

---

## 🛠️ Tech Stack

| Layer         | Technology                       |
| ------------- | -------------------------------- |
| Language      | Golang                           |
| Web Framework | Gin                              |
| DB            | PostgreSQL                       |
| Auth          | JWT, bcrypt                      |
| Logging       | Zap + Lumberjack                 |
| Middleware    | CORS, Custom Logger, RateLimiter |
| Token Store   | Memory (future Redis support)    |

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- (Optional) Mail service for email verify/reset

### Setup

```bash
# clone the repo
git clone https://github.com/thanhdev1710/flamee_auth.git
cd flamee_auth

# config your .env or settings.yaml
cp config/settings.example.yaml config/settings.yaml

# run
go run main.go
```
