<p align="center">
<img src="logo.png" />
</p>
goHOME is an open source home automation client + server, aimed at allowing home owners to have full control over their home automation hardware.

Note - still in alpha development, do not use just yet :)

##Why use an open source home automation project?
###Unified User Interface
I have home automation devices from many different manufacturers in my home, each device requires a different control app, either logging in to a browser app, or installing a mobile phone app, so to control my home I have to use serveral apps. Having one consistent UI where you can control all your devices is a big win for ease of use.  Also, software is not a primary focus of hardware manufacturers, the UIs to control these devices is generally pretty terrible, an after thought. With goHOME the aim is to provide a great user experience and allow the community to customize and provide enahanced experiences not available in the original manufacturer apps.

###Local control - no outside services
One of the main aims of the project is to interop with the hardware on your local network.  Too many home automation devices now require you to connect to an outside service just to talk to the device on your local network.  This is inefficient and error prone.  I have had serveral times where I wanted to change my lights but couldn't because the external service was down, even though the hub on my local network was fine.

###Interoperability
You want all of your devices to be able to work with one another.  By providing an additional unification layer on top of all the hardware devices it is easy for us to plug devices together, without requiring support from the manufacturers.

###Device/Service obsolescence
If you have invested money in buying these hardware devices, it sucks when a manufacturer decides to shut down the service that controls the device or no longer provide support for the control app etc.  By using goHOME the aim is that we can make sure you can control and monitor your devices forever, even long after the manufacturer has given up providing support. 

##Lofty Project Aims
I have a roadmap for where I would like to go with this project:
  - Control Hardware
  - Monitor Hardware
  - Raspberry PI support, allowing goHOME to run 24/7 on a cheap small device
  - Recipes, allow programing locally like IFTTT support
//TODO:

##Supported Hardware
###Lutron Caseta Wireless Smart Bridge
###Flux WIFI Bulbs
###ConnectedByTCP Bulbs
###Belkin WeMo Insight Switch

##Documentation
###Buttons
###Commands
###Devices
###Recipes
###[Scenes](docs/scene.md)
###[Zones](docs/zone.md)
###[HTTP API](docs/api.md)


##Events
 /api/v1/events/ws

##Development
Currently there are two main parts to the project, the golang powered server and the web UI.

###gohome Server
The goHome server is developed using golang (http://golang.org) In order to develop the code:
//TODO:

###gohome web UI
The web UI is developed using the React framework: https://facebook.github.io/react/ In order to develop the web UI:
 1. Setup the goHome Server, following the above instructions
 2. Install node.js: https://nodejs.org
 3. Change to the gohome/www directory
 4. Run:
 
 ```bash
 npm install
 ```
 5. Start the web server and webpack with the --watch option, any file changes are automatically picked up and rebuilt:
 
 ```bash
 npm start
 ```
The server starts on port 8000  http://localhost:8000

####NOTE - All web UI code is located at gohome/www/assets
