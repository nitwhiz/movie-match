FROM node:18.16.0-alpine3.17

RUN npm install --location=global pnpm@8.5.1

COPY docker/dev/app/wrapper.sh /wrapper.sh

WORKDIR /app

USER node

CMD [ "/wrapper.sh" ]
