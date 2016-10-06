<p align="center">
<img src="logo.png" />
</p>
goHOME is an open source home automation client + server, aimed at allowing home owners to have full control over their home automation hardware.

Note - still in alpha development, do not use just yet :)

##Project Aims
Web/iOS/Android/Apple Watch

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
 3. Change to the root gohome directory
 4. Run: "npm install"
 5. Run: "node_modules/webpack/bin/webpack.js --watch --config ./webpack.config.js"
webpack is used to build the React source, the --watch option allows the tool to automatically rebuild the code when it detects any dependant file has changed. This way you just save your modified file and refresh the browser to see the latest changes.
 6. All web UI code is located at gohome/www/assets
