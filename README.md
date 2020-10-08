# KubeCraftAdmin
KubeCraftAdmin : The Adventurer's Admin Tool

![Would you kill this innocent looking service?](https://miro.medium.com/max/700/1*U4MfxStrHa41MUywGgT8ZQ.png)

This project allows you to do basic Kubernetes administration through Minecraft.
[See here](https://medium.com/@eric.jadi/minecraft-as-a-k8s-admin-tool-cf16f890de42) for a more detailed introduction.

Azure Pipeline Status [![Build Status](https://dev.azure.com/ericjadi/KubeCraftAdmin%20-%20Pipelines/_apis/build/status/erjadi.kubecraftadmin?branchName=main)](https://dev.azure.com/ericjadi/KubeCraftAdmin%20-%20Pipelines/_build/latest?definitionId=42&branchName=main)
Latest [Docker Image](https://hub.docker.com/repository/registry-1.docker.io/erjadi/kubecraftadmin)

## Quickstart  

Read this if you're just interested in trying out KubeCraftAdmin yourself.
You will need the following to get started:

- A k8s cluster that you own and are allowed to break
- A place to run the kubecraftcontainer, this can be anywhere inside or outside of the cluster. As long as it has network connectivity to the cluster and to your Minecraft client.
- Minecraft Bedrock Edition

How to run KubeCraftAdmin:

1. Run the container erjadi/kubecraftadmin passing the external port and the location of your .kube directory. The container internally listens to port 8000 and in my case my .kube directory resides in /home/erjadi/.kube
```
 docker run -p 8000:8000 -v /home/erjadi/.kube:/.kube erjadi/kubecraftadmin
```
2. Start up Minecraft Bedrock Edition
3. Create a new world with the *Activate Cheats* option turned **on**
 ![Activate Cheats](/img/cheats.png)
4. Once you've spawned log into the server with the following command after which you should be greeted by the KubeCraftAdmin splash screen.
```
/connect 10.0.0.1:8000/ws
```
5. Next find a nice area to spawn your kubecraft pen. Type *init* to generate the structure.
6. Lastly step on the [beacon](https://minecraft.gamepedia.com/Beacon) to activate the link with your cluster.
7. At this point your pens should be populated with animals!

## Technical Details

### Structure

KubeCraftAdmin is written in Golang. It builds upon the great [MCWSS](https://github.com/Sandertv/mcwss) project by [Sandertv](https://github.com/Sandertv).


### How to compile

The easiest way to compile your own version is to use the provided Dockerfile.
This way you can build it without requiring a local golang build environment.

```
docker build -t kubecraftadmin .
```

