# CloudNativeNow2023

Some commands that you should install
```
brew install protobuf
brew install go
```

Some environment variables to take into consideration
```
export PATH="$PATH:$(go env GOPATH)/bin"
export GO111MODULE=on
```


To install protoc-gen-go
```
go get google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

Install protoc-gen-go-grpc
```
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

Install grpc
To generate the Go and Proto implementation
```



Command to generate Proto and Go code:
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    match.proto
```

Other secret commands to use:
```
go mod tidy
go mod init usac.projects/grpc
```

This links save my life when programming gRPC with Go:
- https://rambabuy.medium.com/understanding-go-go111module-16917777053c
- https://go.dev/doc/tutorial/create-module
- https://developer.confluent.io/get-started/
- https://docs.confluent.io/kafka-clients/go/current/overview.html

## Kafka Strimzi
To install Strimzi run the following commands

```
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl get pod -n kafka --watch
kubectl logs deployment/strimzi-cluster-operator -n kafka -f
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n kafka 
```

To check the installation run
```
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n kafka 
````

To test the installation run
```
kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic

kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
```

To delete Strimzi run
```
kubectl -n kafka delete $(kubectl get strimzi -o name -n kafka)
kubectl -n kafka delete -f 'https://strimzi.io/install/latest?namespace=kafka'
```

If you want to test locally Kafka, just port forward the service and change your /etc/hosts
```
kubectl port-forward svc/my-cluster-kafka-bootstrap 9092:9092 -n kafka
echo "127.0.0.1	localhost my-cluster-kafka-0.my-cluster-kafka-brokers.kafka.svc" >> /etc/hosts
```


# Grafana

To port-forward Grafana run
```
kubectl port-forward -n project svc/grafana 3000:3000
```
To install the plugin take a look into this link:
- https://github.com/RedisGrafana/grafana-redis-datasource


# Locust
To install Locust create a virtual environment using virtualenv library.
To install it run:
```
pip3 install virtualenv
virtualenv env1
source env1/bin/activate
```
for exit to this virtual env run
```
deactivate
```

# Locust
To install locust run
```
pip3 install locust
````
To run locust run the following command:
```
locust
```
Now open your browser in:
```
http://localhost:8089
```

# Kafka topic
To consume pending messages in the matches topic run
```
kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.35.1-kafka-3.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic matches --from-beginning
```
To delete the topic run
```
kubectl delete kafkatopics matches -n kafka
```



# View Logs

Monitor API
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=api -n project --context google) -n project -f api-proj --context google
```

Monitor consumer
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=consumer -n project --context google) -n project -f consumer-proj --context google
```

Monitor server
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=server -n project --context azure) -n project -f server-proj --context azure
```


# Linkerd multi cluster configuration
To install Linkerd and the multicluster configuration run:
```
brew install step
step certificate create root.linkerd.cluster.local root.crt root.key --profile root-ca --no-password --insecure
step certificate create identity.linked.cluster.local issuer.crt issuer.key --profile intermediate-ca --not-after 8760h --no-password --insecure --ca root.crt --ca-key root.key

linkerd install --crds | kubectl --context=google apply -f -
linkerd install --crds | kubectl --context=azure apply -f -

linkerd install --identity-trust-anchors-file root.crt  --identity-issuer-certificate-file issuer.crt --identity-issuer-key-file issuer.key | kubectl --context=google apply -f -
  
linkerd install --identity-trust-anchors-file root.crt --identity-issuer-certificate-file issuer.crt --identity-issuer-key-file issuer.key | kubectl --context=azure apply -f -
  
linkerd --context=google viz install | kubectl --context=google apply -f -
linkerd --context=azure viz install | kubectl --context=azure apply -f -

linkerd --context=google multicluster install | kubectl --context=google apply -f -
linkerd --context=azure multicluster install | kubectl --context=azure apply -f -


linkerd --context=azure multicluster link --cluster-name azure | kubectl --context=google apply -f -
linkerd --context=google multicluster check
linkerd --context=google multicluster gateways


kubectl create deploy webserver --image=nginx --context azure
kubectl get deploy webserver -o yaml --context azure | linkerd inject - | kubectl apply --context google -f -

kubectl expose deploy webserver --target-port=80 --port=80 --context azure
kubectl --context=azure label svc webserver mirror.linkerd.io/exported=true

kubectl create deploy client --image=nginx --context google
kubectl get deploy client -o yaml --context google | linkerd inject - | kubectl apply --context google -f -

kubectl exec -it --context google client-55f7f64b4-9gpqw -c nginx -- sh
curl webserver-azure
```

# How to install this demo
Go to the yaml folder

# Demonstration used commands
Monitor API
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=api -n project --context google) -n project -f api-proj --context google
```

Monitor consumer
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=consumer -n project --context google) -n project -f consumer-proj --context google
```

Monitor server
```
kubectl logs pod/$(kubectl get pods -o jsonpath='{.items[*].metadata.name}' -l app=server -n project --context azure) -n project -f server-proj --context azure
```

kubectl port-forward -n project svc/grafana 3000:3000
Open Grafana in http://localhost:3000

Open Locust http://0.0.0.0:8089 in the browser

Open Azure in Redis and delete the keys to reset the dashboard running
```
del messages
del teams:counters
del teams:phases
```

Open Mongo Compass and drop the collection matches