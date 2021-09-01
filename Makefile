.PHONY:swagger
swagger:
	swag init --output docs/swagger

.PHONY:swagger
baas:
	go build -v -o ./bin/baas

./bin/baas: baas

.PHONY:swagger
start: ./bin/baas
	./bin/baas server