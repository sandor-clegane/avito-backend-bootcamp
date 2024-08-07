generate-mock:
	mockgen -source=./internal/service/flat/interface.go -destination=./internal/service/flat/mocks/mock.go
	mockgen -source=./internal/service/house/interface.go -destination=./internal/service/house/mocks/mock.go
	mockgen -source=./internal/http/handlers/create-flat/handler.go -destination=./internal/http/handlers/create-flat/mocks/mock.go
	mockgen -source=./internal/http/handlers/update-flat/handler.go -destination=./internal/http/handlers/update-flat/mocks/mock.go
	mockgen -source=./internal/http/handlers/get-house/handler.go -destination=./internal/http/handlers/get-house/mocks/mock.go
	mockgen -source=./internal/http/handlers/create-house/handler.go -destination=./internal/http/handlers/create-house/mocks/mock.go
