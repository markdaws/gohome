A feature is a piece of functionality that a device provides.  A feature could be a button, sensor, light zone etc.  Devices can have zero or many features that they export. 

A feature has the following fields:
  - ID
  - Type
  - Address
  - Name
  - Description
  - DeviceID
  - Attrs

//TODO: Attrs

Below is a list of all the of the different types of features that goHOME currently supports.

##Button
##Cool Zone
##Heat Zone
A heat zone represents an output from your furnace.  Your furnace might support multiple zones, meaning different parts of your house can be set to different temperatures, or you might just have one zone where the whole house is the same temperature. A heat zone supports the following attributes:
###Current Temperature
###Target Temperature

##Light Zone
A light zone can be thought of as one or more bulbs that are all controlled at the same time.  Think of it as a piece of wire with one or more bulbs attached to it.  A light zone supports the following attributes:
###On/Off
###Brightness (optional)
###Hue (optional)
All of the bulbs in a zone are set to the same values, you can't control them individually.  Light zones generally map to the physical wiring in your house.

##Outlet
##Sensor
##Switch
##Window Treatment
A Window Treatment represents some mechanism to cover your windows, either shades, curtains or something else. A window treatment supports the following attributes
