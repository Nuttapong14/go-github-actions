# Frontend Dockerfile - Multi-stage build for Next.js application
# Dependencies stage
FROM oven/bun:1-alpine AS deps

WORKDIR /app

# Copy package files
COPY frontend/package.json frontend/bun.lockb* ./

# Install dependencies
RUN bun install --frozen-lockfile || bun install

# Builder stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY frontend/ ./

# Set Next.js telemetry disabled
ENV NEXT_TELEMETRY_DISABLED=1

# Build the application
RUN bun run build

# Production stage
FROM oven/bun:1-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set environment variables
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Copy necessary files from builder
COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next && chown appuser:appgroup .next

# Copy the standalone build and static files
COPY --from=builder --chown=appuser:appgroup /app/.next/standalone ./
COPY --from=builder --chown=appuser:appgroup /app/.next/static ./.next/static

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# Run the application
CMD ["bun", "run", "server.js"]
