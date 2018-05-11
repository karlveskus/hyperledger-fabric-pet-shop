# Hyperledger Fabric Pet Shop

## Requirements
- MacOS / Linux / VirtualBox (required if Windows)
- [Git](https://git-scm.com/downloads)
- [Node](https://nodejs.org/en/)
- [Go language](https://golang.org/dl/)
- [Docker](https://www.docker.com/)

## Setup
1. Run ```./bootstrap.sh``` to get hyperledger fabric binaries
2. ```cd pet-shop```
3. Remove pre-existing docker containers ```docker rm -f $(docker ps -aq)```
4. Run ```./startFabric.sh``` to setup network
5. Run ```npm install``` to get node modules
6. Run ```node registerAdmin.js && node registerUser.js``` to create admin and user
7. Run ```node server.js``` to start client server
8. Open ```http://localhost:8000```


## Troubleshooting

### Error running ```startFabric.sh```
```
ERROR: failed to register layer: rename
/var/lib/docker/image/overlay2/layerdb/tmp/write-set-091347846 /var/lib/docker/image/overlay2/layerdb/sha256/9d3227c1793b7494e598caafd0a5013900e17dcdf1d7bdd31d39c82be04fcf28: file exists
```

try running ```rm -rf ~/Library/Containers/com.docker.docker/Data/*```

### Error running ```node registerAdmin.js && node registerUser.js```
```
Error: [client-utils.js]: sendPeersProposal - Promise is rejected: Error: Connect Failed

error from query =  { Error: Connect Failed

   at /Desktop/prj/education/LFS171x/fabric-material/tuna-app/node_modules/grpc/src/node/src/client.js:554:15 code: 14, metadata: Metadata { _internal_repr: {} } }
```
or
```
Failed to register: Error: fabric-ca request register failed with errors [[{"code":20,"message":"Authorization failure"}]]
Authorization failures may be caused by having admin credentials from a previous CA instance.
```

try running ```rm -rf ~/.hfc-key-store/``` and then registering command again.

## Finishing up
To remove all Docker containers and images run following commands
1. ```docker rm -f $(docker ps -aq)```
2. ```docker rmi -f $(docker images -a -q)```