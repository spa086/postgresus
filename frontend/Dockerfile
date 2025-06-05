FROM node:18-alpine

EXPOSE 4006

WORKDIR /app

COPY package*.json ./
RUN npm ci && npm install -g serve

COPY . .

RUN npm run build

CMD ["serve", "-s", "dist", "-l", "4006"]
