default:
	GOOS=linux go build -v -o example
	docker build -t pkger:example .
	docker run -p 3000:3000 pkger:example
