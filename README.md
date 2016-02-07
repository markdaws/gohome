# gohome
golang powered home automation server

Note - this is not ready for production yet, it's running all the lights in my house but I am still working on getting a v1 that supports saving/loading and setup via the UI vs. hacking JSON files.  Should be ready for use in the next few weeks (mid March), if you try to run any of this code before then, you're going to have a bad time.

#Project Aims
Web/iOS/Android/Apple Watch

#Supported Hardware
###Lutron Caseta Wireless Smart Bridge
###Flux WIFI Bulbs
###ConnectedByTCP Bulbs
###Belkin WeMo Insight Switch
Ideally we would support as many different types of hardware as possible, if you want to see support for a piece of hardware submit an issue.

#Core Concepts

##Devices
##Zones
A Zone is a controllable output. For example a zone could be a piece of wire with several bulbs attached to it, or it could be a shade, or an outlet. It's something that you can control.  In the case of lights, generally a Zone will be the wire with light attached to it that then connects to a switch or dimmer. Currently there are several types of Zones:
 - Light
 - Shade
 - Outlet

The type lets the UI know what controls to show for a particular Zone.  Along with the Type a Zone has an Output value which can be:
 - Continuous: Something that can be set between 0 and 100 e.g. a dimmable bulb
 - Binary: Something that can be set on or off e.g. an outlet
 - RGB: A Zone that doesn't have just one value, but can display different colors e.g. Flux WIFI bulb

##Buttons
##Commands
###Supported Commands
#####ZoneSetLevel
#####SceneSet
#####ButtonPress
#####ButtonRelease
##Scenes
A Scene is simply a group of commands that will be executed when you activate the Scene. Examples of Scenes could be:
- "All On" -> sets all of the lights in your house to 100%
- "Movie Time" -> turns off all the lights in your living room, closes all of your shades

Since a scene is just a group of commands you can mak a scene do literally anything.  A Scene can activate other scenes, control lights, control shades, anything that goHome has command support for.  See the Commands section for more information on the list of supported commands.

###Recipes/Actions/Triggers

##API Support
###Scenes
GET /api/v1/scenes
GET /api/v1/scenes/<ID>
DELETE /api/v1/scenes/<ID>
PUT /api/v1/scenes/<ID>
POST /api/v1/scenes
POST /api/v1/scenes/active
POST /api/v1/scenes/<ID>/commands
DELETE /api/v1/scenes/<ID>/commands/<INDEX>

###Zones
GET /api/v1/zones
POST /api/v1/zones
PUT /api/v1/zones/<ID>

###Discovery
GET /api/v1/discovery/<MODEL_NUMBER>
GET /api/v1/discovery/<MODEL_NUMBER>/token
GET /api/v1/discovery/<MODEL_NUMBER>/access
GET /api/v1/discovery/<MODEL_NUMBER>/zones

###Events
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
