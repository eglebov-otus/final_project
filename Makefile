test:
	@docker build --target unit-test .

lint:
	@docker build --target lint .

build:
	@docker build --target build .

up-server:
	@docker build --target server -t image-previewer:latest .
	@docker run -d -p 8080:8080 --name image-previewer image-previewer:latest

down-server:
	@docker container stop image-previewer
	@docker container rm image-previewer
	@docker image rm image-previewer:latest

stop-server:
	@docker container stop image-previewer

start-server:
	@docker container start image-previewer