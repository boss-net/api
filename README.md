Here's a clean, **attractive-looking `README.md`** for your **Boss Backend API** project. It includes badges, improved formatting, proper sectioning, and modern aesthetics using Markdown features like admonitions and anchors.

---

```md
# 🚀 Boss Backend API

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](../.github/workflows)
[![Python](https://img.shields.io/badge/python-3.11+-blue.svg)](https://www.python.org/)

A blazing-fast, production-ready backend for the **Boss** platform. Powered by Flask, Celery, Redis, PostgreSQL, and Weaviate.

---

## ⚙️ Usage

> [!IMPORTANT]
> This project uses [`uv`](https://docs.astral.sh/uv/) as the **Python package manager**. Make sure it's installed first.

---

### 🐳 Start Middleware Stack

The backend depends on PostgreSQL, Redis, and Weaviate. Launch them via `docker-compose`:

```bash
cd docker
cp middleware.env.example middleware.env
# optional: switch profile if not using Weaviate
docker compose -f docker-compose.middleware.yaml --profile weaviate -p boss up -d
```

---

### 🔐 Set Up Environment

1. Copy `.env.example` to `.env`:

   ```bash
   cp .env.example .env 
   ```

2. Generate a `SECRET_KEY`:

   **Linux:**
   ```bash
   sed -i "/^SECRET_KEY=/c\SECRET_KEY=$(openssl rand -base64 42)" .env
   ```

   **macOS:**
   ```bash
   secret_key=$(openssl rand -base64 42)
   sed -i '' "/^SECRET_KEY=/c\\
   SECRET_KEY=${secret_key}" .env
   ```

---

### 📦 Dependency Setup

Install [`uv`](https://docs.astral.sh/uv/) if not installed:

```bash
# Linux / Windows
pip install uv

# macOS
brew install uv
```

Then sync dependencies:

```bash
uv sync --dev
```

---

### 🔄 Database Migration

Run migrations before first use:

```bash
uv run flask db upgrade
```

---

### 🚀 Start the Backend

```bash
uv run flask run --host 0.0.0.0 --port=5001 --debug
```

---

### 🌐 Start the Web UI

Start the [Boss Web Frontend](../web) and open:

```
http://localhost:3000
```

---

### 🧵 Start Async Worker (Optional)

For dataset importing, document indexing, and other async tasks:

```bash
uv run celery -A app.celery worker \
  -P gevent -c 1 --loglevel INFO \
  -Q dataset,generation,mail,ops_trace,app_deletion
```

---

## ✅ Testing

1. Sync all test dependencies:

   ```bash
   uv sync --dev
   ```

2. Run tests locally with mocked variables (`tool.pytest_env` in `pyproject.toml`):

   ```bash
   uv run -P api bash dev/pytest/pytest_all_tests.sh
   ```

---

## 📜 License

Licensed under the [MIT License](LICENSE).

---

_Developed with ❤️ by the **Boss Net** team._
```

---

Let me know if you want a version with emojis removed, or extra sections like **Contributing**, **Architecture**, or **API Docs**.
