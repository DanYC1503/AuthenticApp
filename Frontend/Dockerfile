# Stage 1: Build Angular app
FROM node:18 AS build
WORKDIR /app
COPY auth_frontend/ .
RUN npm install
RUN npm run build -- --configuration production

# Stage 2: Serve with Nginx
FROM nginx:alpine
COPY --from=build /app/dist/auth_frontend /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
