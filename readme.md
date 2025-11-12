# 🧩 OMHS Backend

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-ready-blue?logo=docker)](https://www.docker.com/)
[![MongoDB](https://img.shields.io/badge/MongoDB-Atlas%20%2F%20Local-green?logo=mongodb)](https://www.mongodb.com/)

> Backend service for the **OMHS Project**, built with Go, Gin, and MongoDB.  
> Provides authentication, data management, and integration endpoints for the OMHS platform.

---

## 📖 Overview

The **OMHS backend** exposes REST APIs for authentication, users, and requests.  
It’s written in **Go**, uses **MongoDB** for storage, and leverages **Gin** as the HTTP router and middleware layer.  
Docker Compose is used to manage both the backend and the database for local and production setups.

---

## ⚙️ Tech Stack

| Category | Technology |
|-----------|-------------|
| Language | [Go 1.25+](https://golang.org) |
| Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | [MongoDB](https://www.mongodb.com/) |
| Auth | Session tokens (non‑JWT) |
| Containerization | [Docker](https://www.docker.com/) |

---

## 🚀 Quick Start

### 🧩 Development

Run the backend and MongoDB locally:

```bash
docker compose -f docker-compose.dev.yml up --build
```

Stop containers:

```bash
docker compose -f docker-compose.dev.yml down
```

Access backend at:  
👉 [http://localhost:8080](http://localhost:8080)

---

### 🏗️ Production

Build and run the production image:

```bash
docker compose up --build -d
```

Stop services:

```bash
docker compose down
```

---

## 🔐 Environment Variables

Your `.env` file (ignored by Git) should include:

```env
MONGO_URI=mongodb://mongo:27017
PORT=8080
TOKEN_EXPIRATION_HOURS=48
```

> 📝 **Tip:** Make sure `.env` is included in `.gitignore`.

---

## 🧰 Common Commands

| Action | Command |
|--------|----------|
| Run dev environment | `docker compose -f docker-compose.dev.yml up --build` |
| Stop dev containers | `docker compose -f docker-compose.dev.yml down` |
| Run production | `docker compose up --build -d` |
| Stop production | `docker compose down` |
| View logs | `docker compose logs -f backend` |

---

## 💡 Notes

- Backend runs on **port 8080**
- MongoDB data is persisted via Docker volume
- CORS is configured for frontend origin `http://localhost:4200`
- Session tokens are stored and validated via MongoDB (non‑JWT approach)

---

## 🧑‍💻 Author

**Omri Meller**  
🕒 Last updated: **November 2025**
