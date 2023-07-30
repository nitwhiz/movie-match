FROM node:18.16.0-alpine3.17 AS base

ENV PNPM_HOME="/app/.pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

COPY docker/dev/app/wrapper.sh /wrapper.sh

WORKDIR /app

CMD [ "/wrapper.sh" ]
