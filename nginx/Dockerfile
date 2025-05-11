FROM nginx:alpine

# SSL 証明書をコピー
COPY ssl/server.crt /etc/nginx/ssl/server.crt
COPY ssl/server.key /etc/nginx/ssl/server.key

# nginx 設定ファイルをコピー
COPY nginx.conf /etc/nginx/nginx.conf

# ログディレクトリを作成
RUN mkdir -p /var/log/nginx

# ポートを公開
EXPOSE 80 443

# nginx を起動
CMD ["nginx", "-g", "daemon off;"]
