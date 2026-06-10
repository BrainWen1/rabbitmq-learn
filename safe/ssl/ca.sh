#!/bin/bash
# 创建根证书颁发机构（CA）
openssl genrsa -out ca_key.pem 2048
openssl req -x509 -new -nodes -key ca_key.pem -days 3650 -out ca_certificate.pem -subj "/CN=MyCA"

# 创建服务端私钥
openssl genrsa -out server_key.pem 2048

# 1. 创建配置文件 ssl.conf（关键：添加 SAN 扩展，避免浏览器/客户端报错）
cat > ssl.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name

[req_distinguished_name]

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
IP.1 = 192.168.12.143  # 当前虚拟机IP
IP.2 = 127.0.0.1 # 本地回环地址
DNS.1 = rabbitmq-server  # 也可以改成主机名，方便域名访问
DNS.2 = localhost
EOF

# 2. 生成服务器证书请求（包含 SAN 扩展）
openssl req -new -key server_key.pem -out server.csr -subj "/CN=rabbitmq-server" -config ssl.conf

# 3. 使用 CA 签名服务器证书（包含 SAN 扩展）
openssl x509 -req -in server.csr \
  -CA ca_certificate.pem \
  -CAkey ca_key.pem \
  -CAcreateserial \
  -out server_certificate.pem \
  -days 365 \
  -extensions v3_req \
  -extfile ssl.conf

# 生成客户端私钥和证书（用于客户端连接验证）
openssl genrsa -out client_key.pem 2048
openssl req -new -key client_key.pem -out client.csr -subj "/CN=rabbitmq-client"
openssl x509 -req -in client.csr -CA ca_certificate.pem -CAkey ca_key.pem -CAcreateserial -out client_certificate.pem