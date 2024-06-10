LOCATION=westus
RESOURCE_GROUP=vtd-rg-$RANDOM
SERVICE_BUS_NAMESPACE=vtd-sb-$RANDOM
SERVICE_BUS_TOPIC_LOCATION=locations
SERVICE_BUS_TOPIC_SHIPMENTS=shipments
COSMOSDB_ACCOUNT=vtd-cosmos-$RANDOM
COSMOSDB_COLLECTION=vehicle-tracker
COSMOSDB_DATABASE=vehicle-tracker
AKS_CLUSTER=vtd-aks-$RANDOM

# Create a resource group
az group create --name $RESOURCE_GROUP --location $LOCATION

# create an AKS cluster
az aks create --resource-group vtd-rg-24601 --name $AKS_CLUSTER --node-count 1 --enable-addons monitoring --generate-ssh-keys

# create a service bus namespace
az servicebus namespace create --resource-group $RESOURCE_GROUP --name $SERVICE_BUS_NAMESPACE --location $LOCATION

# create a topic for locations
az servicebus topic create --resource-group $RESOURCE_GROUP --namespace-name $SERVICE_BUS_NAMESPACE --name $SERVICE_BUS_TOPIC_LOCATION

# create a topic for shipments
az servicebus topic create --resource-group $RESOURCE_GROUP --namespace-name $SERVICE_BUS_NAMESPACE --name $SERVICE_BUS_TOPIC_SHIPMENTS

# get the connection string
az servicebus namespace authorization-rule keys list --resource-group $RESOURCE_GROUP --namespace-name $SERVICE_BUS_NAMESPACE --name RootManageSharedAccessKey --query primaryConnectionString --output tsv

#save the connection string in a variable
$SERVICE_BUS_CONNECTION_STRING = az servicebus namespace authorization-rule keys list --resource-group $RESOURCE_GROUP --namespace-name $SERVICE_BUS_NAMESPACE --name RootManageSharedAccessKey --query primaryConnectionString --output tsv
kubectl create secret generic sb-connection-string --from-literal=sb-connection-string=$SERVICE_BUS_CONNECTION_STRING -n vehicle-tracker

# create a cosmos db account
az cosmosdb create --name $COSMOSDB_ACCOUNT --resource-group $RESOURCE_GROUP --kind GlobalDocumentDB --locations regionName=$LOCATION failoverPriority=0 isZoneRedundant=False

# get the connection string
#az cosmosdb list-connection-strings --name $COSMOSDB_ACCOUNT --resource-group $RESOURCE_GROUP --query connectionStrings[0].connectionString --output tsv

#save the connection string in a variable
$COSMOSDB_CONNECTION_STRING = az cosmosdb list-connection-strings --name $COSMOSDB_ACCOUNT --resource-group $RESOURCE_GROUP --query connectionStrings[0].connectionString --output tsv
kubectl create secret generic cosmosdb-connection-string --from-literal=cosmosdb-connection-string=$COSMOSDB_CONNECTION_STRING -n vehicle-tracker

# create a database
az cosmosdb database create --name $COSMOSDB_ACCOUNT --db-name $COSMOSDB_DATABASE --resource-group $RESOURCE_GROUP

# create a collection
az cosmosdb collection create --collection-name $COSMOSDB_COLLECTION --name $COSMOSDB_ACCOUNT --db-name $COSMOSDB_DATABASE --resource-group $RESOURCE_GROUP --partition-key-path /partitionKey --throughput 400
