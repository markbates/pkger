default:
	cd ../../../cmd/pkger && go install -v -tags debug .
	pkger
	GOOS=linux go build -v -tags debug -o example
	docker build -t pkger:example .
	docker run -p 3000:3000 pkger:example
