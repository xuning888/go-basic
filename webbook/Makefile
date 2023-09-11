.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build --tags=k8s -o webook .
	@docker rmi -f xuning888/webook:v.0.0.1
	@docker build -t xuning888/webook:v.0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go --package=svcmock -destination=./internal/service/mocks/user.mock.go
	@mockgen -source=./internal/service/code.go --package=svcmock -destination=./internal/service/mocks/code.mock.go
	@go mod tidy