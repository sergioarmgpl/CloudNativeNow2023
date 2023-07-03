# CloudNativeNow2023

export PATH="$PATH:$(go env GOPATH)/bin"
export GO111MODULE=on
go mod init grpc-football-match
brew install protobuf
brew install go

Install protoc-gen-go

$ go get google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go

Install protoc-gen-go-grpc

$ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc 
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

Install grpc

$ go get google.golang.org/grpc
mkdir match;
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    match/match.proto 


protoc --go-grpc_out=require_unimplemented_servers=false:./match/ --go_out=./match/ match.proto

go mod tidy


https://go.dev/doc/tutorial/create-module



go mod init usac.projects/match
https://rambabuy.medium.com/understanding-go-go111module-16917777053c
https://go.dev/doc/tutorial/create-module
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    match.proto
go mod tidy

go mod init usac.projects/grpc



## Kafka
kubectl create namespace project
kubectl create namespace kafka

kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl get pod -n kafka --watch
kubectl logs deployment/strimzi-cluster-operator -n kafka -f
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n kafka 
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n kafka 

kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic


kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning

kubectl -n kafka delete $(kubectl get strimzi -o name -n kafka)

kubectl -n kafka delete -f 'https://strimzi.io/install/latest?namespace=kafka'

kubectl port-forward svc/my-cluster-kafka-bootstrap 9092:9092 -n kafka

127.0.0.1	localhost my-cluster-kafka-0.my-cluster-kafka-brokers.kafka.svc


https://developer.confluent.io/get-started/
https://docs.confluent.io/kafka-clients/go/current/overview.html


kubectl logs deploy/api -n project -f

host=34.170.66.53
curl -X POST -H "Content-Type: application/json" --data '{"team1" : "Korea","team2":"Italy","score":"0-0","phase":16}' http://$host:3000/match


kubectl port-forward -n project svc/grafana 3000:3000

pip3 install locust

locust
http://localhost:8089
http://34.31.92.137:3000


kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic matches --from-beginning

kubectl delete kafkatopics matches -n kafka

https://github.com/RedisGrafana/grafana-redis-datasource


# View Logs
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=api -n project --context google) -n project -f api-proj --context google

kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=consumer -n project --context google) -n project -f consumer-proj --context google


kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=server -n project --context azure) -n project -f server-proj --context azure




kubectl port-forward -n project svc/grafana 3000:3000 --context google