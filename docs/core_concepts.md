In order to better understand the UI and potentially write automation, it is essential to understand a few core concepts.

##Devices
A device is a piece or hardware that you can monitor or control. For example a device might be a wall dimmer, or a WIFI Bulb, or a WIFI Outlet, or a window shade remote control. In your home automation system you will have one or more devices that you want to monitor and control together. goHOME lets you put all these devices in one common UI, saving you from having to use multiple apps, one from each hardware manufacturer.

The important properties for a device are:
  - ID: a unique ID that allows the device to be identified
  - Name: a name that describes the device
  - Address: This could be an IP address or some other identifier, that allows us to connect and control the device
  - Features: See below.

##Features
A device has one or more features. A feature is a piece of functionality, for example, the Belkin WeMo Maker device is a piece of hardware that has a switch you can open/close and a sensor that detects an open/close state. In this example we have one device and two features, a switch and a sensor.  A light dimming device would export a light zone, a thermostat could export a heat zone for monitor/control temperature and potentially a fan if your furnace supports it.

The currently supported types of feature are:
###Button
A physical or virtual button that a person may press. For exampel a light switch device might export a on/off button.

###Cool Zone/Heat Zone
A heat zone represents an output from your furnace.  Your furnace might support multiple zones, meaning different parts of your house can be set to different temperatures, or you might just have one zone where the whole house is the same temperature.

###Light Zone
A light zone can be thought of as one or more bulbs that are all controlled at the same time.  Think of it as a piece of wire with one or more bulbs attached to it. All of the bulbs in a zone are set to the same values, you can't control them individually.  Light zones generally map to the physical wiring in your house.

###Outlet
###Sensor
A sensor can monitor any kind of property, open/closed, temperature, light level etc.

###Switch

###Window Treatment
A Window Treatment represents some mechanism to cover your windows, either shades, curtains or something else. 

##Scenes
A Scene in simplest terms can be though of as a collection of commands. For example, you might create a "Movie" scene which you activate when you are watching a movie at home. The scene will:

  - Close the living room shades
  - Set the living room lights to 10% intensity
  - Turn on the popcorn maker (via a Belkin WeMo Switch)

Other examples of scenes might be "All lights off", "Relaxing", "Dinner Time". You can specify a list of commands that will be executed sequentially when the scene is activated.

##Extensions
Extensions allow goHOME to be extended to support different kinds of hardware. To read more about extensions and how to create them, see <a href="extensions.md">here</a>
