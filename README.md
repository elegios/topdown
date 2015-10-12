# Topdown
A client-server structure for top down turnbased games. It features a strong separation between client and server, exstensibility through scripting on the server and a roguelike feel. The protocol is specified in [protocol.txt](protocol.txt).

## Design Ideas

The goal is to provide a structure for a long running game with a mix of MMO elements and co-op story driven games. The world should always be running and provide something to do, as well as larger events that happen once (generally when most of the playrs on the server are online). It should also be easy to extend the world and the story by adding more worlds and functionality to those worlds.

## Extension Structure

The idea is to have a bunch of modules, each containing maps, map updates and stories. Maps contain all geography, map updates can be applied by stories to change the look of the world whenever the story calls for it. Stories are scripts written in lua that either run in the background, controlling elements in the world, or runs once to set something else up.

Examples of this can be found in [world/modules](world/modules) in the ```main``` module (which is loaded when the server first starts up) and the ```ancient_temple``` module, which is loaded by the ```main``` module.

When restarting the server remembers which modules were applied, which map updates (called partials) were applied and where, and which stories were running. Each story has a key-value store for data that needs to be remembered through a reboot, and the main function of each story gets a single parameter, a bool, telling the story whether it should do first time setup or recovery.
