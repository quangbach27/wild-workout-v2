# Stage 1: Builder
FROM node:24-alpine3.24 AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Vite inlines VITE_* variables at build time.
ARG VITE_API_URL
ARG VITE_MOCK_USER
ARG VITE_FUTURE_WEEKS
RUN npm run build

# Stage 2: Runtime
FROM nginx:stable-alpine3.24-slim

COPY docker/app-prod/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/app-prod/security-headers.conf /etc/nginx/snippets/security-headers.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
