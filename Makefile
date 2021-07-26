.PHONY:swagger
swagger:
	swag init --output docs/swagger

baas:
	go build -v -o ./bin/baas

./bin/baas: baas

start: ./bin/baas
	./bin/baas server