# Vehicle Tracker Dapr

A sample shipment tracker to showcase Dapr components and features.

## Architecture and Service Definitions



## Demo setup

## Prerequisites

- Kubernetes cluster of your choice (3+ nodes recommended)
- Helm

## Clone this repository

```git
git clone https://github.com/rochabr/vehicle-tracker-dapr.git

cd vehicle-tracker-dapr
```

## Create application namespaces

After connecting to your cluster, run the following command to create the namespaces:

```bash
# create application namespaces
kubectl create ns vehicle-tracker # vehicle teracker namespace
kubectl create ns dapr-system # Dapr
kubectl create ns redis # Redis
kubectl create ns kafka # Kafka
kubectl create ns zipkin # Zipkin
```

## Install Dapr and create RBAC

You can choose to skip this step if you want Conductor to manage your Dapr installation (recommended):

```bash
dapr init -k -n dapr-system
kubectl apply -f ./components/minikube/dapr-secret-reader.yaml
```

## Setup Helm

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

## Redis setup

Install Redis, export the password and create a secret that will be accessed from the component files.

```bash
helm install redis bitnami/redis -n redis
export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 -d) 

kubectl create secret generic redis-password --from-literal=redis-password=$REDIS_PASSWORD -n vehicle-tracker
```

## Kafka setup

Install Kafka, export the password and create a secret that will be accessed from the component files.

```bash
helm install --set persistence.enabled=false --set zookeeper.persistence.enabled=false --set auth.clientProtocol=sasl kafka bitnami/kafka -n kafka

export KAFKA_PASSWORD=$(kubectl get secret kafka-user-passwords --namespace kafka -o jsonpath='{.data.client-passwords}' | base64 -d | cut -d , -f 1)
kubectl create secret generic kafka-password --from-literal=kafka-password=$KAFKA_PASSWORD -n vehicle-tracker
```

## Install Zipkin

Zipkin will be used for applicaiton tracing.

```bash
echo "Zipkin namespace created"
kubectl create deployment zipkin -n zipkin --image openzipkin/zipkin
echo "Zipkin deployment created"
kubectl expose deployment zipkin -n zipkin --type LoadBalancer --port 9411 
echo "Zipkin deployment exposed"
kubectl get svc -n zipkin -w
export ZIPKIN_DASHBOARD=$(kubectl get svc --namespace zipkin zipkin -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9411
echo "View tracing dashboard at $ZIPKIN_DASHBOARD"
```

## Deploy Dapr components

Now we will deploy the components that we will use throughout the demo:

```bash
kubectl apply -f ./components/minikube
```

## Deploy services

```bash
kubectl apply -f ./deployment/services
```

> The contents of the folder `/deployment-files/k8s` are used for automated CI/CD pipelines, as described in the main [README](./../README.md) file.

Check the state of your pods by running:

```bash
kubectl get pods -n vehicle-tracker
```

## Expose the service on MiniKube

```bash
minikube service shipment-handler --url -n vehicle-tracker
```

