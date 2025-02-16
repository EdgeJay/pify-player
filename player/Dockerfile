FROM node:22 AS builder

WORKDIR /app

COPY . .
RUN npm install
RUN npm run build
RUN npm ci --omit dev

FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app/build/ /app/build/

ENV PORT=3000

EXPOSE 3000

CMD ["node", "./build"]