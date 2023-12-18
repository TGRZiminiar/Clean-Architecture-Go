## Template For Clean Architecture Golang
Design for microservice starter template
- using grpc to communicate between microservice
- using postgres as database

## Command Generate Private Key Pem
```
openssl genpkey -algorithm RSA -out private-key.pem
cat private-key.pem | openssl pkey -text
```