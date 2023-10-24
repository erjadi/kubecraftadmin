# KubeCraftAdmin
KubeCraftAdmin : The Adventurer's Admin Tool

![Would you kill this innocent looking service?](https://miro.medium.com/max/700/1*U4MfxStrHa41MUywGgT8ZQ.png)

This project allows you to do basic Kubernetes administration through Minecraft.
- [See here](https://medium.com/@eric.jadi/minecraft-as-a-k8s-admin-tool-cf16f890de42) for a more detailed introduction.
- Link to the latest [Docker Image](https://hub.docker.com/repository/registry-1.docker.io/erjadi/kubecraftadmin)
- Latest build status ![Build Status](https://dev.azure.com/ericjadi/KubeCraftAdmin%20-%20Pipelines/_apis/build/status/erjadi.kubecraftadmin?branchName=main)

## Quickstart  

Read this if you're just interested in trying out KubeCraftAdmin yourself.
You will need the following to get started:

- A k8s cluster that you own and are allowed to break
- A place to run the kubecraftcontainer, this can be anywhere inside or outside of the cluster. As long as it has network connectivity to the cluster and to your Minecraft client.
- Minecraft Bedrock Edition

How to run KubeCraftAdmin:

1. Run the container erjadi/kubecraftadmin passing the external port and the location of your .kube directory. The container internally listens to port 8000 and in my case my .kube directory resides in /home/erjadi/.kube. You can optionally specify one to four namespaces from your cluster using the environment variable *namespaces*
```
 docker run -p 8000:8000 -v /home/erjadi/.kube:/.kube [-e namespaces=mynamespace1,mynamespace2] erjadi/kubecraftadmin
```
2. Start up Minecraft Bedrock Edition
3. Create a new world with the *Activate Cheats* option turned **on**
 ![Activate Cheats](/img/cheats.png)
4. Once you've spawned log into the server with the following command after which you should be greeted by the KubeCraftAdmin splash screen. You should also be given some items to start with (a sword, TNT and some flint).
```
/connect 10.0.0.1:8000/ws
```
5. Next find a nice area to spawn your kubecraft pen. Type *init* to generate the structure.
6. Lastly step on the [beacon](https://minecraft.wiki/w/Beacon) to activate the link with your cluster.
7. At this point your pens should be populated with animals!

## Running KUbecraftadmin on OpenShift

- oc new-project kubecraft
- oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:kubecraft:default
- oc create -f deploy.yaml
- Get the route
- Load minecraft and use ```/connect <route-url>/ws```

## Technical Details

### Structure

KubeCraftAdmin is written in Golang. It builds upon the great [MCWSS](https://github.com/Sandertv/mcwss) project by [Sandertv](https://github.com/Sandertv).

This project makes use of the [Websocket Server](https://minecraft.wiki/w/Commands/wsserver) functionality present in Minecraft Bedrock and Education Edition. The WS connection is a Minecraft client connection, which means that all actions are performed through the client. The server / local world is unaffected and not controlled by this project. This also implies we need to activate cheats in the world to be able to [summon](https://minecraft.wiki/w/Commands/summon) or [kill](https://minecraft.wiki/w/Commands/kill) entities.

The below description explains the main process which you can find in [kubecraftadmin.go](/src/app/kubecraftadmin.go).  
Highly simplified, KubeCraftAdmin connects to the Kubernetes cluster, spawns the required entities and starts an endless loop function *LoopReconcile*. Every second it starts a sync function called *ReconcileKubetoMC* which basically:

- Enumerates entities in Minecraft
- Enumerates resources in Kubernetes
- Kills / Spawns the differences in Minecraft

For the reverse sync we rely on a mobEvent which triggers execution of *ReconcileMCtoKubeMob*.
We basically perform the same check, but this time we take the Minecraft entities as the truth and delete the corresponding Kubernetes resources.

### Known Issues / TODO

- There are some hacks to make the sync stable. Some operations take time on the kubernetes cluster (e.g. you kill a service in Minecraft). For the duration of the deletion process, Minecraft and Kubernetes will be out of sync (chicken is dead, service is still there). This could lead to syncing issues where KubeCraftAdmin tries to spawn chickens. Right now this is 'fixed' by keeping a list of uniqueIDs and **never spawning the same entity twice**. This has its own issues. A slight improvement would be to have a TTL on the uniqueIDs, but the best solution would be to keep track of the state as "being synced to MC" or "being synced to K8s".  
- Hardcoded to 4 namespaces (more than 4 will crash)
- Logging is bad / non-configurable
- No error handling

### How to compile

The easiest way to compile your own version is to use the provided Dockerfile.
This way you can build it without requiring a local golang build environment.

```
docker build -t kubecraftadmin .
```

