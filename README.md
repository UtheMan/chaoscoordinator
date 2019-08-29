# Chaos Coordinator
Chaos Coordinator is a set of tools that allow for chaos testing of the infrastructure used by Kubernetes clusters on Azure. 

Kubernetes cron jobs are used to ensure periodic testing that follows user-defined schedule.
 
# Motivation
Chaos Coordinator aims to make chaos testing of kubernetes clusters infrastructure as streamlined as possible.
Chaos tests are triggered by Kubernetes cron jobs, which frees the user from the need of keeping track of the current state of tests and schedules.
This makes it possible for the user to focus solely on the type of chaos user might want to implement.

For more information on chaos testing please refer to the [Principles of Chaos Engineering](https://principlesofchaos.org/?lang=ENcontent).
 
# Technologies used
* Go
* Kubernetes
* Docker
* [azure-sdk-for-go](https://github.com/Azure/azure-sdk-for-go)
* [client-go](https://github.com/kubernetes/client-go)
* [Skaffold](https://github.com/GoogleContainerTools/skaffold)
* Chaos CLI - [cobra](https://github.com/spf13/cobra)
# Features
* Creates stateless chaos testing scenarios using REST API calls.
* Schedules of created scenarios are managed by Kubernetes cron jobs.
* Allows for extension - implement new types of chaos using Go, as needed for your specific use case.
* Deployable on your Kubernetes cluster.
* Usable locally as a binary.
 
# How to use Chaos Coordinator
To correctly use Chaos Coordinator you need a Kubernetes cluster running on Azure. Rest of this readme assumes you already have a cluster up and running.
## Structure
Chaos Coordinator is made of 2 parts:
* [Chaos CLI](https://github.com/UtheMan/chaoscoordinator/tree/master/cmd)
* [Chaos Coordinator API](https://github.com/UtheMan/chaoscoordinator/tree/master/cron)
### Chaos Coordinator CLI
Chaos Coordinator CLI implements the CLI which is used by cron jobs to trigger chaos. Available commands can be seen [here](https://github.com/UtheMan/chaoscoordinator/blob/master/cmd/chaos.go).
All the details regarding specific commands can be found in the [/cmd](https://github.com/UtheMan/chaoscoordinator/tree/master/cmd) package.
 
Example commands:
* Create a cron job that reboots a random VM in a scale set in a given resource group
```
./chaos vm kill -m random -n SCALE_SET_NAME -r RESOURCE_GROUP_NAME
```
* Create a cron job that adds 1GB of data to the disk space of a given VMSS for 120 seconds  
```
./chaos disk fill -d 120 -a 1000 -n SCALE_SET_NAME -r RESOURCE_GROUP_NAME
```
### Chaos Coordinator API
The API is an entry point of Chaos Coordinator, allows the user to create, delete and get details of cron jobs responsible for chaos testing.   
API is available under ```API_SERVICE_ADDRESS/api```

Example requests:
* POST - create a cron job called test which reboots a VM in VMSS every minute
```
URL: http://API_SERVICE_ADDRESS/api
Request body:
{
    "name": "test",
    "schedule": "*/1 * * * *",
    "command": ["./chaos"],
    "args": ["vm", "kill", "-m", "random", "-r", "YOUR_RESOURCE_GROUP_NAME", "-n", "YOUR_SCALE_SET_NAME"]
}
```
* GET - get information about cron job with name "test"
```
URL: http://API_SERVICE_ADDRESS/api?name=test
```
* GET - get information about all cron jobs in "default" namespace
```
URL: http://API_SERVICE_ADDRESS/api?namespace=default
```
* DELETE - delete cron job with name "test"
```
URL: http://API_SERVICE_ADDRESS/api?name=test
```
## Secrets
To correctly authenticate with Azure, ChaosCoordinator requires two secrets to be present on kubernetes cluster:
* SUBSCRIPTION_ID - your Azure subscription ID.     
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: azure-subscription-id
type: Opaque
data:
  subscriptionId: YOUR_AZURE_SUBSCRIPTION_ID
```
Create with ```kubectl apply -f filename.yaml```
 
* AZURE_AUTH_LOCATION - your Azure credentials file,
```yaml
{
  "clientId": YOUR_CLIENT_ID,
  "clientSecret": YOUR_CLIENT_SECRET,
  "subscriptionId": YOUR_SUBSCRIPTION_ID,
  "tenantId": YOUR_TENANT_ID,
  "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
  "resourceManagerEndpointUrl": "https://management.azure.com/",
  "activeDirectoryGraphResourceId": "https://graph.windows.net/",
  "sqlManagementEndpointUrl": "https://management.core.windows.net:8443/",
  "galleryEndpointUrl": "https://gallery.azure.com/",
  "managementEndpointUrl": "https://management.core.windows.net/"
}
```
Create with ```kubectl create secret generic azure-auth --fromfile=creds=filename```
## Build
Docker Images are built and deployed with Skaffold. For more information please refer to [Skaffold documentation](https://skaffold.dev/docs/getting-started/#installing-skaffold).

Both API and CLI can be built and launched locally. Use ```make``` command in respective directories to trigger appropriate builds.   
 
## Deploy
To use Chaos Coordinator on your cluster you need to:
* Create a Kubernetes service on your cluster to route traffic to Chaos Coordinator API pods
```yaml
apiVersion: v1
kind: Service
metadata:
  name: chaoscoordinatorsvc
  labels:
    run: chaoscoordinatorservice
spec:
  ports:
  - port: 80
    targetPort: 3000
    protocol: TCP
  type: LoadBalancer
  selector:
    run: chaoscoordinatorservice
```
Add to cluster with ```kubectl apply -f SERVICE_FILE_NAME.yaml```
* Add Chaos Coordinator API deployment to your cluster. 
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chaoscoordinatorservice
spec:
  selector:
    matchLabels:
      run: chaoscoordinatorservice
  replicas: 2
  template:
    metadata:
      labels:
        run: chaoscoordinatorservice
    spec:
      containers:
        - name: chaoscoordinatorservice
          image: utheman/utheman_chaoscoordinatorservice:latest
          command: [./bin/chaoscoordinatorservice]
          ports:
          - containerPort: 3000
```
Add to cluster with ```kubectl apply -f DEPLOYMENT_FILE_NAME.YAML```
## Creating Chaos
TODO
## New commands for the CLI
Chaos Coordinator CLI can be extended with additional commands as needed, 
Cobra allows for easy extension - see [here](https://github.com/spf13/cobra) for more information regarding Cobra.  
For an example of a command used in Chaos Coordinator see [here](https://github.com/UtheMan/chaoscoordinator/tree/master/cmd/vm).
### Adding new commands
All the commands for Chaos Coordinator CLI reside in the /cmd package.

To register new commands with the CLI one has to add a subcommand to the [root command](https://github.com/UtheMan/chaoscoordinator/blob/master/cmd/chaos.go).
Follow vm command as an [example](https://github.com/UtheMan/chaoscoordinator/tree/master/cmd/vm) while structuring your commands.
### Command implementation
Every command in the CLI calls it's implementation as seen [here](https://github.com/UtheMan/chaoscoordinator/blob/7322e51ade8bc5f2e96b9550160d829dd956d2b8/cmd/vm/kill/kill_vm.go#L20).
These implementations are located in the [pkg/cmd/azure](https://github.com/UtheMan/chaoscoordinator/tree/master/pkg/cmd/azure) package.

This code is later executed by the Kubernetes cron jobs deployed on the cluster. 
# License
This project is licensed under the [MIT License](https://github.com/UtheMan/chaoscoordinator/blob/master/LICENSE).

