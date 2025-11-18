# 🧩 OMHS Backend

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://golang.org)  
[![Docker](https://img.shields.io/badge/Docker-ready-blue?logo=docker)](https://www.docker.com/)  
[![MongoDB](https://img.shields.io/badge/MongoDB-Atlas%20%2F%20Local-green?logo=mongodb)](https://www.mongodb.com/)

Backend service for the **OMHS Project**, built with Go, Gin, and MongoDB.  
Provides authentication, data management, and integration endpoints for the OMHS platform.

---

## 📖 Overview

The **OMHS backend** exposes REST APIs for authentication, users, and requests.  
It’s written in **Go**, uses **MongoDB** for storage, and leverages **Gin** as the HTTP router and middleware layer.  
Docker Compose (via Makefile) is used to manage backend + database for both dev and production.

---

## ⚙️ Tech Stack

| Category | Technology |
|----------|-------------|
| Language | Go 1.25+ |
| Framework | Gin |
| Database | MongoDB |
| Auth | Session tokens (non-JWT) |
| Containerization | Docker & Docker Compose |

---

### 🧩 Requirements
- Go 1.25+
- Docker & Docker Compose
- Make  
  - macOS/Linux: already installed  
  - Windows:

    ```powershell
    winget install Ezwinports.Make
    ```

---

### ▶️ Development

Start backend + MongoDB:

```bash
make omhs-dev
```

Stop containers:

```bash
make omhs-dev-down
```

Logs:

```bash
make omhs-dev-logs
```

---

### 🏗 Production

```bash
make omhs-prod
```

---

## 🔧 Most Common Commands

```bash
make omhs-dev       # Start dev environment
make omhs-dev-down  # Stop dev environment
make omhs-dev-logs  # View logs
make omhs-run       # Run Go app locally
make omhs-build     # Build Go binary
```

For the full list of available commands:

```bash
make omhs-help
```

---

## 🔐 Environment Variables

Create a `.env` file:

```env
MONGO_URI=mongodb://mongo:27017
PORT=8080
TOKEN_EXPIRATION_HOURS=48
```

---

## 💡 Notes

- Backend runs on **http://localhost:8080**
- MongoDB data is persisted via Docker volume
- CORS allows `http://localhost:4200`
- Session tokens are validated using MongoDB (non-JWT)

---

## 🧑‍💻 Author

**Omri Meller**  
🕒 Last updated: **November 2025**
