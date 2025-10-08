# Event Management REST API

A simple REST API built with **Go (Gin)** for managing events, users, and registrations.  
Includes **JWT authentication**, **Swagger documentation**, and **Docker** support.

---

## 🧰 Requirements

- Go ≥ **1.25**
- `make`, `docker`, and optionally `air` (for hot reload)

---

## 🚀 Quick Start

### Development

```bash
make tools      # install swag & air
make swag       # generate Swagger docs
make dev        # run with Air (hot reload)
```

App runs on **http://localhost:8089** (Swagger at `/swagger/index.html`).

### Production

```bash
make build           # compile binary
make docker-build    # build Docker image
make docker-run      # run container
```

---

## ⚙️ Environment

Create `.env` file in the project root:

```
JWT_SECRET=change-me
PORT=8080
```

Docker automatically mounts `.env` into the container.

---

## 📚 Swagger Docs

- Path: `/swagger/index.html`
- Command: `make swag`
- Security: `Bearer <JWT>`

---

## 📦 Structure

```
rest-api/
├─ main.go
├─ db/
├─ models/
├─ routes/
├─ middlewares/
├─ utils/
├─ docs/
├─ Makefile
├─ Dockerfile
└─ .env.example
```

---

## 🧠 Common Commands

| Command             | Description             |
| ------------------- | ----------------------- |
| `make dev`          | Run app with hot reload |
| `make run`          | Run app normally        |
| `make swag`         | Generate Swagger docs   |
| `make build`        | Build binary            |
| `make docker-build` | Build Docker image      |
| `make docker-run`   | Run container           |
| `make clean`        | Remove artifacts        |

---

## 🗝 Endpoints Overview

| Route                       | Method      | Auth | Description             |
| --------------------------- | ----------- | ---- | ----------------------- |
| `/signup`                   | POST        | —    | Create new user         |
| `/login`                    | POST        | —    | Sign in and get JWT     |
| `/events`                   | GET         | —    | List all events         |
| `/events`                   | POST        | ✅   | Create event            |
| `/events/:id`               | GET         | —    | Get event by ID         |
| `/events/:id`               | PUT         | ✅   | Update event            |
| `/events/:id`               | DELETE      | ✅   | Delete event            |
| `/events/:id/registrations` | POST/DELETE | ✅   | Register / cancel event |

---

## 🧹 Cleanup

```bash
make clean
```
