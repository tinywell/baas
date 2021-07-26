.PHONY:swagger
swagger:
	swag init --output docs/swagger

baas:
	go build -v

start: baas
	./baas server