prepare:
	@docker build --target base .

test:
	@docker build --target unit-test .

test-race:
	@docker build --target unit-test-race .

lint:
	@docker build --target lint .

build:
	@docker build --target build .

run:
	@docker build --target server -t image-previewer:latest .
	@docker run -d -p 8080:8080 --name image-previewer image-previewer:latest

logs:
	@docker logs -f image-previewer

stop:
	@docker container stop image-previewer
	@docker container rm image-previewer
	@docker image rm image-previewer:latest