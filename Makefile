## Run Test Suite inside docker
tests:
	@docker-compose -f docker-compose.yaml up --build --abort-on-container-exit
	@docker-compose -f docker-compose.yaml down --volumes

## Populate influxdb with some test data and run tests
## Note: This command is intended to be executed within docker env
integration-tests:
	@until nc -z -v -w30 influxdb 8086; do echo Waiting for database connection...; sleep 3; done
	go test -v ./...
