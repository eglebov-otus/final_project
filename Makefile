prepare:
	@docker build --target base .

test:
	@docker build --target unit-test .

lint:
	@docker build --target lint .

build:
	@docker build --target build .

build-server:
	@docker build --target server -t image-previewer:latest .
	@docker run -d -p 8080:8080 --name image-previewer image-previewer:latest

remove-server:
	@docker container stop image-previewer
	@docker container rm image-previewer
	@docker image rm image-previewer:latest