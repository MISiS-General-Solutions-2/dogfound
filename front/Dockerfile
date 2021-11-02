FROM node:17-alpine3.12
WORKDIR /app

ENV NODE_OPTIONS=--openssl-legacy-provider

COPY package.json yarn.lock ./
RUN cd /app \
    && yarn install --pure-lockfile

COPY . /app
CMD ["yarn","start"]