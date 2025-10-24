FROM postgres:15-alpine

# Set environment variables
ENV POSTGRES_DB=auth_system
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin123

# Copy initialization scripts - they'll execute in alphabetical order
COPY 01_schema.sql /docker-entrypoint-initdb.d/
COPY 02_data.sql /docker-entrypoint-initdb.d/

# Expose PostgreSQL port
EXPOSE 5432

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pg_isready -U admin -d auth_system