
build-dev:
	VERSION="dev" goreleaser --snapshot --skip-publish --rm-dist


build-prod:
	VERSION=$(shell cat version) goreleaser --snapshot --skip-publish --rm-dist

install-prod:
	$(eval VER=$(shell cat version))
	go install -ldflags="-X 'local.only/gmcat/cmd.Version=$(VER)'"


docker-dev: build-dev
	$(eval VER=$(shell cat version))
	docker build -t gmcat-dev:$(VER) -f docker/Dockerfile . --no-cache

docker-prod: build-prod
	$(eval VER=$(shell cat version))
	docker build -t gmcat:$(VER) -f docker/Dockerfile . --no-cache