The aim of the goHOME project is to allow users to monitor and control all of their home automation hardware from one UI.  In order to understand how to use the app and how to write extensions etc. there are some core concepts you need to understand.

##Devices
A device is a piece of hardware. A device could be a dimmer, shade, garage door opener, sensor etc. Devices then export features.

##Features
A feature is a piece of functionality the user will monitor and control. For example a feature might be a switch that you want to be able to turn on/off or a HeatZone which you want to be able to set the temperate for. Devices can export many features. For example the Belkin WeMo Maker device has a switch that you can turn on/off and also a sensor that can detect and open/close state, so in goHOME terminology, we would have a device that represents the actual WeMo Maker device then export two features from it, a switch and a sensor.

The current supported features are:

  - Button
  - CoolZone
  - HeatZone
  - LightZone
  - Outlet
  - Sensor
  - Switch
  - WindowTreatment: something that covers a window e.g. a shade, curtain. You want to be able to control and monitor the offset of the window treatment
