# rm *.pem

# 1. Generate CA's private key and self-signed certificate
# openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=BR/ST=Sao Paulo/L=Sao Paulo/O=Casa/OU=Tecnologia/CN=*.tec.com/emailAddress=silas@tec.com"

# echo "CA's self-signed certificate"
# openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-cert.pem -subj "/C=BR/ST=Sao Paulo/L=Sao Paulo/O=Office/OU=Computer/CN=*.tec1.com/emailAddress=silas1@tec1.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in client-cert.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

echo "Server's signed certificate"
openssl x509 -in client-cert.pem -noout -text