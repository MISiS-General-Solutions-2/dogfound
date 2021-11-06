FROM node:17-alpine3.12 as build
WORKDIR /app

ENV NODE_OPTIONS=--openssl-legacy-provider

COPY package.json ./
COPY yarn.lock ./
RUN yarn install --production

# set to your ip
ENV REACT_APP_API_URL="ip:port"

COPY . /app
CMD ["yarn","start"]