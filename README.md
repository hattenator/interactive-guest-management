A client and server for negotiating settings between an interactive qemu guest and its host.

This repository contains two independent Go directories:

* **guest** –  
* **host**  –  

Both modules are fully self‑contained; you can build and run them
separately, or import the packages from the other module if you wish
to share code.

This relies on a channel, which you can add with the provided libvirt/channel.xml.
It should be named after the domain (VM name), so you may need to edit channel.xml (my domain is called "win11").  Then attach it with:
sudo virsh attach-device $domain ./libvirt/channel.xml --config --live

The channel will add instantly, but Windows will likely have some trouble with it until you reboot Windows.
