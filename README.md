# Distributed Final Project - Garuda (Twitter Clone)

## Pre-requisites

To start this project you need to have go installed with all the mnodules mentioned in the project.

You will also need to install etcd and goreman to utilize the raft storage option in the project. 

## Protobuf Generator

Run the ./gen-model.sh shell script found in the root directory to generate the stub codes from the proto files. 

## Starting Application

execute the following commands to start the application

```bash
go run web/web.go
```

The above command will start the webserver which will render our html pages from the locally stored go template files (gtpl) using handler functions written in the same module. 

```bash
go run cmd/web/web.go
```
The above command will start the backend server which will communicate to the webserver through grpc using protobuf files found in the model directory.   

```bash
goreman -f Procfile start
```
The above command will start an etcd cluster of three nodes with endpoints mentioned in the Procfile.

Note - Make sure to edit the procfile to reflect the correct path to the etcd binary. 

