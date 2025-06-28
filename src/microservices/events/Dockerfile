FROM node:18-alpine
WORKDIR /app

RUN apk add --no-cache --virtual .gyp python3 make g++

COPY package*.json ./
RUN npm ci --production

COPY . .
EXPOSE 8082
CMD ["node","index.js"]
