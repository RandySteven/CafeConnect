gorun = @go run

cmd_folder = ./cmd/main/
yaml_file = ./files/yml/cofeConnect.local.yaml

run:
	${gorun} ${cmd_folder}http -config ${yaml_file}

migration:
	${gorun} ${cmd_folder}migration -config ${yaml_file}

seed:
	${gorun} ${cmd_folder}seed -config ${yaml_file}

drop:
	${gorun} ${cmd_folder}drop -config ${yaml_file}

alter:
	${gorun} ${cmd_folder}alter -config ${yaml_file}

consumer:
	${gorun} ${cmd_folder}consumer -config ${yaml_file}

topic:
	${gorun} ${cmd_folder}topic -config ${yaml_file}

server: scheduler run consumer

build:
	@go build -o /bin/http

refresh: drop migration seed

scheduler:
	${gorun} ${cmd_folder}scheduler -config ${yaml_file}

docker_run:
	docker compose build && docker compose run library-app

test:
	go test -coverprofile=coverage.out ./... ;    go tool cover -html=coverage.out

mockery:
	mockery --all --output ./mocks

start_nsq:
	nsqd \\n  --lookupd-tcp-address=127.0.0.1:4160 \\n  --broadcast-address=127.0.0.1 \\n  --tcp-address=127.0.0.1:4150 \\n  --http-address=127.0.0.1:4151\n
	nsqlookupd
	nsqadmin --lookupd-http-address=127.0.0.1:4161