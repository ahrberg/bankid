# use this script to convert p12 file to pem files
# password: qwerty123

openssl pkcs12 -in test_client_cert.p12 \
    -clcerts -nokeys -out client_cert.pem
openssl pkcs12 -in test_client_cert.p12 \
    -nocerts -out client_key.pem -nodes