.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f xuning888/webook:v.0.0.1
	@docker build -t xuning888/webook:v.0.0.1 .