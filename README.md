Here’s your fixed and cleaned-up `README.md` version:

---

# AuthenticApp - Secure Authentication System

A **microservices-based authentication system** with an Angular frontend, built with **Go** and **PostgreSQL**.

---

## 🚀 Quick Start

### Prerequisites

* Docker
* Docker Compose

### Local Development

1. Clone this repository
2. Run:

   ```bash
   docker-compose up -d --build
   ```
3. Access:

   * Frontend → [http://localhost:4200](http://localhost:4200)
   * API Gateway → [http://localhost:8888](http://localhost:8888)

---

## 🧩 Services

| Service           | Port | Description                                     |
| ----------------- | ---- | ----------------------------------------------- |
| **API Gateway**   | 8888 | Routes requests, handles CORS & CSRF protection |
| **Auth Service**  | 9999 | JWT authentication & user management            |
| **User Service**  | 8889 | User data operations                            |
| **Audit Service** | 8890 | Activity logging                                |
| **PostgreSQL**    | 5432 | Data storage                                    |

---

## 👤 Default Admin

* **Username:** `admin123`
* **Email:** `admin@gmail.com`

---

## 🏗 Architecture

```
Angular Frontend → API Gateway → Auth Service / User Service / Audit Service → PostgreSQL
```

---

## 🔧 Deployment

### Production

```bash
docker-compose -f docker-compose.prod.yml up -d
```

---

Would you like me to include a small diagram (ASCII or image-based) of the architecture flow? It could make the README look more professional.
