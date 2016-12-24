<p align="center">
<img src="logo.png" />
</p>
goHOME is an open source home automation project, aimed at allowing home owners to have full control over their home automation hardware in a single UI. The project is designed to run on cheap hardware like a Raspberry PI.

Here it is in action, controlling some Lutron Lights and a Belkin WeMo Insight Outlet
<p align="center">
<img src="https://github.com/markdaws/gohome-assets/blob/master/gohome_demo_720.gif" />
</p>

##Why use an open source home automation project?
###Unified User Interface
I have home automation devices from many different manufacturers in my home, each device requires a different control app, either logging in to a browser app, or installing a mobile phone app, so to control my home I have to use serveral apps. Having one consistent UI where you can control all your devices is a big win for ease of use.  Also, software is not a primary focus of hardware manufacturers, the UIs to control these devices is generally pretty terrible, an after thought. With goHOME the aim is to provide a great user experience and allow the community to customize and provide enahanced experiences not available in the original manufacturer apps.

###Local control - no outside services
One of the main aims of the project is to interop with the hardware on your local network.  Too many home automation devices now require you to connect to an outside service just to talk to the device on your local network.  This is inefficient and error prone.  I have had serveral times where I wanted to change my lights but couldn't because the external service was down, even though the hub on my local network was fine.

###Interoperability
You want all of your devices to be able to work with one another.  By providing an additional unification layer on top of all the hardware devices it is easy for us to plug devices together, without requiring support from the manufacturers.

###Device/Service obsolescence
If you have invested money in buying these hardware devices, it sucks when a manufacturer decides to shut down the service that controls the device or no longer provide support for the control app etc.  By using goHOME the aim is that we can make sure you can control and monitor your devices forever, even long after the manufacturer has given up providing support. 

##Supported Hardware
For a list of the supported hardware, click <a href="docs/supported_hardware.md">here</a>.  If you want support for a piece of hardware that is not on the list, file an issue.

##Documentation
###[Core Concepts](docs/core_concepts.md)
###[Getting Started](docs/getting_started.md)
###[Supported Hardware](docs/supported_hardware.md)
###[Raspberry PI Setup](docs/raspberrypi_manual.md)
###[FAQ](docs/faq.md)
###[Automation](docs/automation.md)
###[Development](docs/development.md)

##Feedback
If you have any issues, you can file a Github issue on this project, if you have any questions, feel free to email me: mark@gohome.io
