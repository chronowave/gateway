docker-gateway: build-web
	docker build -t chronowave/gateway .

build-web:
	yarn install
	yarn build

.PHONY: docker-gateway build-web
