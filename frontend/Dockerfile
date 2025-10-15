# build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json tsconfig*.json ./
COPY src ./src
RUN npm ci --silent
RUN npm run build

# final stage - serve with a small Node static server
FROM node:20-alpine AS runner
WORKDIR /app
COPY --from=builder /app/dist ./dist
RUN npm install -g serve --silent
EXPOSE 80
CMD ["serve", "-s", "dist", "-l", "80"]
