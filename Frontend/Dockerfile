# ===============================
# Stage 1: Build the Angular app
# ===============================
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files first for caching
COPY Auth_App_Frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the Angular source code
COPY Auth_App_Frontend/ .

# Build the Angular app for production
RUN npm run build -- --configuration production

# =================================
# Stage 2: Serve the app with Nginx
# =================================
FROM nginx:alpine

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy Angular build output to Nginx web root
COPY --from=builder /app/dist/auth-app-frontend/. /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Healthcheck (optional)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --spider -q http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
