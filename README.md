# kubecraftadmin
KubeCraftAdmin : The Adventurer's Admin Tool

This project allows you to do basic Kubernetes administration through Minecraft.
[See here](https://medium.com/@eric.jadi/minecraft-as-a-k8s-admin-tool-cf16f890de42) for a more detailed introduction.

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
4. Once you've spawned log into the server with the following command
```
/connect 10.0.0.1:8000/ws
```
This should show the kubecraft splash screen
5. Next find a nice area to spawn your kubecraft pen. Type *init* to generate the structure.
6. Lastly step on the [beacon](https://minecraft.gamepedia.com/Beacon) to activate the link with your cluster.
7. At this point your pens should be populated with animals!

# Technical Details

TODO

