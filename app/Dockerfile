FROM node:18.14-slim AS builder

RUN npm i -g pnpm@7.27.0

WORKDIR /build

COPY ./ /build

RUN pnpm install && pnpm build

FROM nginx:1.23.3-alpine

COPY --from=builder /build/dist/ /usr/share/nginx/html/
