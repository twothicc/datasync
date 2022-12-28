
all: vendors format lint compile

fieldAlignment:
	fieldalignment -fix github.com/twothicc/datasync/config

format:
	gofmt -s -w $$(find . -type f -name '*.go'| grep -v "/vendor/")

lint:
	golangci-lint run

vendors:
	go mod vendor

compile:
	cd app; go build -o ../build/datasync

clearLog:
	> server.log

# Connect to jaeger via http://localhost:16686
startJaegerUI:
	docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.9

# Connect to prometheus via http://localhost:9090
startPrometheus:
	docker run -p 9090:9090 --name prometheus -v C:/GoApps/datasync/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

startElasticKibana:
	cd docker-compose; docker-compose --env-file ../.env up

start:
	./build/datasync