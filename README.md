# SortLink - Modern URL Shortener & Realtime Analytics

A fullstack URL Shortener built with **Golang (Gin, GORM)**, **PostgreSQL (Docker)**, and **Vue 3 + Shadcn Vue + Tailwind CSS**.

---

## ✨ Features

- **🔐 Multi-User Authentication**: JWT-based user authentication (Register, Login, Me).
- **⚡ Fast Redirection**: High-throughput URL resolution with asynchronous analytics recording.
- **🏷️ Custom Slug / Alias**: Create custom branded short URLs (e.g. `localhost:8080/promo-special`).
- **🔒 Password Protection**: Secure links with a passcode before redirecting visitors.
- **⏳ Expiration Date**: Set an auto-expiration date for time-limited campaigns.
- **🎯 Click Limit**: Automatically deactivate links after reaching a target click count.
- **📊 Realtime Analytics & Dashboard**:
  - Time-series click graphs (7, 14, 30 days).
  - Breakdown by Device (Mobile, Desktop, Tablet), Operating System, and Browser.
  - Traffic Referrer sources & Geolocation (Country and City).
- **📥 CSV Data Export**: Download raw click event data for reports or analysis.
- **📱 QR Code Generator**: Generate and download high-resolution QR codes (PNG).
- **🌗 Dark Mode & Light Mode**: Theme switching with persistent user preference.

---

## 🛠️ Tech Stack

- **Backend**: Golang 1.22+, Gin Web Framework, GORM ORM, Golang-JWT, Bcrypt, Mileusna UserAgent.
- **Database**: PostgreSQL 16 (via Docker Compose).
- **Frontend**: Vue 3 (Composition API, `<script setup>`, TypeScript), Vite, Tailwind CSS v4, Shadcn Vue, Pinia, Vue Router, Chart.js, QRCode, Lucide Icons.

---

## 🚀 Getting Started

### Quick Start with Makefile

```bash
# 1. Setup environment and install dependencies
make setup

# 2. Run backend & frontend concurrently (with DB auto-started)
make dev
```

### Individual Service Commands

| Command | Description |
|---|---|
| `make help` | Tampilkan daftar semua perintah makefile |
| `make setup` | Copy `.env`, install frontend dependencies & jalankan PostgreSQL |
| `make dev` | Jalankan Backend dan Frontend bersamaan |
| `make db-up` / `make db-down` | Start / Stop container PostgreSQL |
| `make db-logs` | Lihat logs database PostgreSQL |
| `make db-reset` | Reset database dan volume |
| `make be-run` | Jalankan backend Go server (`http://localhost:8080`) |
| `make be-build` | Build binary backend ke `backend/bin/server` |
| `make be-test` | Jalankan unit tests Go backend |
| `make be-tidy` | Jalankan `go mod tidy` |
| `make fe-install` | Install node modules frontend |
| `make fe-dev` | Jalankan frontend Vite dev server (`http://localhost:5173`) |
| `make fe-build` | Build frontend untuk production |
| `make fe-typecheck` | Jalankan type check TypeScript frontend |
| `make clean` | Hapus build binary dan file dist |

---

## 🧪 Testing

### Backend Unit Tests
```bash
make be-test
# atau: cd backend && go test -v ./...
```

### Frontend Build Verification
```bash
make fe-build
# atau: cd frontend && npm run build
```
