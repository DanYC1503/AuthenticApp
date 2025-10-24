# AuthenticApp - Secure Authentication System

A microservices-based authentication system with Angular frontend, built with Go and PostgreSQL.

## 🚀 Quick Start

### Prerequisites
- Docker
- Docker Compose

### Local Development
1. Clone this repository
2. Run: `docker-compose up -d --build`
3. Access:
   - Frontend: http://localhost:4200
   - API: http://localhost:8888

### Services
- **API Gateway**: 8888 - Request routing & CSRF protection
- **Auth Service**: 9999 - JWT authentication & user management  
- **User Service**: 8889 - User data operations
- **Audit Service**: 8890 - Activity logging
- **PostgreSQL**: 5432 - Data storage

### Default Admin
- Username: `admin123`
- Email: `admin@gmail.com`

## 🏗 Architecture
