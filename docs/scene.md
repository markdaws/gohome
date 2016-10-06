A Scene in simplest terms can be though of as a collection of commands.  For example, you might create a "Movie" scene which you activate when you are watching a movie at home.  The scene will:
  1. Close the living room shades
  2. Set the living room lights to 10% intensity
  
Other examples of scenes might be "All lights off", "Relaxing", "Dinner Time".  You can specify a list of commands that will be executed sequentially when the scene is activated.

##Managed vs Unmanaged Scenes
An unmanaged scene is a scene that is controlled by a 3rd party device.  For example you may import a scene from a Lutron lighting device into the goHOME system.  In this case goHOME has no idea what commands will run when the scene is executed, so you can't modify the behaviour of this scene in any way from within the goHOME UI. 

On the other hand, a managed Scene is a Scene that goHOME knows which commands are being executed and can be edited in the goHOME UI.

##Scene Commands
The following commands can be executed by a scene:

###Button Press
Presses a button on a device

###Button Release
Releases a button on a device

###Scene Set
Sets the specified scene

###Zone Set Level
Sets the level on the specified zone

###Zone Turn On
Turns the zone on to full intensity

###Zone Turn Off
Turns the zone off. Note some devices can use set level == 0 to turn off, others require an explicit off command.
