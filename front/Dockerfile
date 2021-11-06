FROM node:17-alpine3.12 as build
WORKDIR /app

ENV NODE_OPTIONS=--openssl-legacy-provider

COPY package.json ./
COPY yarn.lock ./
RUN yarn install --production

COPY . /app
CMD ["yarn","start"]