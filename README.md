# AuthenticApp - Sistema de Autenticación Segura

Un **sistema de autenticación basado en microservicios** con un frontend en **Angular**, desarrollado con **Go** y **PostgreSQL**.

---

## 🚀 Inicio Rápido

### Prerrequisitos

* Docker
* Docker Compose

### Desarrollo Local

1. Clona este repositorio
2. Ejecuta:

   ```bash
   docker-compose up -d --build
   ```
3. Accede a:

   * Frontend → [http://localhost:4200](http://localhost:4200)
   * API Gateway → [http://localhost:8888](http://localhost:8888)

---

## 🧩 Servicios

| Servicio          | Puerto | Descripción                                       |
| ----------------- | ------ | ------------------------------------------------- |
| **API Gateway**   | 8888   | Enruta solicitudes, maneja CORS y protección CSRF |
| **Auth Service**  | 9999   | Autenticación JWT y gestión de usuarios           |
| **User Service**  | 8889   | Operaciones sobre datos de usuario                |
| **Audit Service** | 8890   | Registro de actividades                           |
| **PostgreSQL**    | 5432   | Almacenamiento de datos                           |

---

## 👤 Administrador Predeterminado

* **Usuario:** `admin123`
* **Correo:** `admin@gmail.com`

---

## 🏗 Arquitectura

```
Angular Frontend → API Gateway → Auth Service / User Service / Audit Service → PostgreSQL
```

---

## 🔧 Despliegue

### Producción

```bash
docker-compose -f docker-compose.prod.yml up -d
```

---

### Architecture.md

```markdown
# Arquitectura del Sistema

## Stack Tecnológico
- **Frontend**: Angular con Tailwind CSS  
- **Backend**: Microservicios en Go  
- **API Gateway**: Go con Gorilla Mux  
- **Base de Datos**: PostgreSQL con claves primarias UUID  
- **Autenticación**: Tokens JWT con protección CSRF  
- **Contenerización**: Docker y Docker Compose  

## Seguridad
- Hashing de contraseñas con sal  
- Autenticación basada en tokens JWT  
- Protección CSRF  
- Registro de auditorías para todas las acciones  
- Validación de entradas y prevención de inyecciones SQL  

## Decisiones de Escalabilidad
- Arquitectura de microservicios para escalar de forma independiente  
- API Gateway para enrutamiento centralizado  
- Autenticación sin estado (stateless)  
- Despliegue contenerizado con Docker  
```
