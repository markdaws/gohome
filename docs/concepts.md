The aim of the goHOME project is to allow users to monitor and control all of their home automation hardware from one UI.  In order to understand how to use the app and how to write extensions there are some core concepts you should familiarize yourself with.

##Devices
A device is a piece of hardware. A device could be a dimmer, shade, garage door opener, sensor etc. It is a piece of hardware you can monitor or control. Devices also know how to convert goHOME commands like "turn light on" into hardware specific commands.

##Features
A feature is a piece of functionality the user will monitor and control. For example, a feature might be a switch that you want to be able to turn on/off or a HeatZone which you want to be able to set the temperate for. Devices can export many features. For example the Belkin WeMo Maker device has a switch that you can turn on/off and also a sensor that can detect and open/close state, so in goHOME terminology, we would have a device that represents the actual WeMo Maker device then export two features from it, a switch and a sensor.

The current supported features are:

  - Button - either a physical or virtual button that users can press
  - CoolZone - a zone that can receive cool air from your home cooling system
  - HeatZone - a zone that can receive heated air from you home heating system
  - LightZone - one or more bulbs that can be controlled at the same time
  - Outlet - a high voltage outlet
  - Sensor - a generic sensir. Sensors can monitor any property, the sensor describes the type of data it is monitoring
  - Switch - a on/off switch
  - WindowTreatment: something that covers a window e.g. a shade, curtain. You want to be able to control and monitor the offset of the window treatment

##Attributes
A feature can have one or more attributes.  Attributes are just values, for example the HeatZone feature has a "current temperature" and "target temperature" attribute.  You can query a feature for all of the latest values for its attributes and you can set values in the attributes.

##Scenes
A scene can be thought of as one or more actions that should be executed when the scene is activated. For example, you might have a "movie" scene, this scene executes the following commands:
  - Close shades in living room
  - Turn off all lights in living room
  - Turn on the outlet which has the popcorn maker plugged in to it
  
Scenes can be as complicated as you like, scenes can also execute other scenes if desired.

##Users
goHOME supports multiple users. At present there are no user specific customizations available, but they will be added soon, such as per user ordering of elements in the UI, hiding UI elements, styling UI etc. Even though there are presently no user specific customizations, it may be useful to have different users so that you can see who logged in to the system and the actions they performed, which are stored in the events.json log.

##Areas
All devices and features live in a specific area.  An area is just a physical space, such as "bathroom", "kitchen" etc.  Presently all devices and features live in the "Home" area.  Areas are not yet exposed in the UI but are present in the data model and will be exposed in the UI shortly.
