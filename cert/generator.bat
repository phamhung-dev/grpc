:: delete pem files
del /f *.pem

:: create ca private key and self-signed certificate
:: adding -nodes to not encrypt the private key
openssl req -x509 -newkey rsa:4096 -nodes -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=VN/L=Ho Chi Minh/O=Pham Hung/CN=*.phamhung-dev/OU=Training"

:: create web server private key and CSR (Certificate Signing Request)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=VN/L=Ho Chi Minh/O=Pham Hung/CN=*.phamhung-dev/OU=Server"

:: sign the web server certificate request (CSR)
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf

:: verify certificate
openssl verify -CAfile ca-cert.pem server-cert.pem

:: create client's private key and CSR (Certificate Signing Request)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=VN/L=Ho Chi Minh/O=Pham Hung/CN=*.someclient.com/OU=Client"

:: sign the client certificate request (CSR)
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf

:: verify certificate
openssl verify -CAfile ca-cert.pem client-cert.pem