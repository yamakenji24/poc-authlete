FROM node:18-alpine

WORKDIR /app

# ホスト名解決のために必要なパッケージをインストール
RUN apk add --no-cache curl

COPY package.json package-lock.json ./
RUN npm install

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]
