# Stage 1: Dependencies and build
FROM node:20-alpine AS builder
WORKDIR /app

# Copy only package files first (better caching)
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# Stage 2: Run production server
FROM node:20-alpine
WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy built app and set ownership
COPY --from=builder --chown=appuser:appgroup /app ./

USER appuser
EXPOSE 3000
CMD ["npm", "start"]
