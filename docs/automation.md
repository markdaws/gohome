//TODO: Write config doc

Automation allows you to program your home devices to execute commands based on certain conditions. For example, you might want to turn on the outside lights at sunset or open all your window shades at sunrise.  You also might want to respond to other events such as turning on all your lights if the garage door sensor changes to an 'open' state after 10pm.

In order to write automation there are a few basic concepts you need to understand. 

##Triggers & Actions
Automation consists fo two main parts:
 - Trigger: a trigger is an event or set of conditions that specify when you want your automation to execute. For example, a trigger may be a time trigger, that will execute at a specific date and time, or at 4pm every day or at sunrise. A trigger may also be an event like 'sensor state changed' or 'button pressed' which you want to be able to respond to.
 
 - Actions: actions are all of the actions that you wish to execute when the trigger fires. For example you may specify multiple actions to execute when the "sunrise" trigger fires, such as: "turn on kitchen lights", "open bedroom shades", "turn on coffee maker"
 
##Creating Automation Scripts
Currenly there is no UI to create automation scripts, they are written by hand.
 
###File Location
When writing automation scripts, you need to tell goHOME the directory where all your scripts are located. By default goHOME will look in the directory where the gohome executable is located, under a directory call "automation".  If you want to change this location you can modify your config.json file, see here //TODO: for more details on modifying the config file.

Each automation file can contain one piece of automation, so a single trigger that fires one or more actions.  If you want multiple pieces of automation e.g. turn outside lights on at 7pm and another to turn on sprinklers at 5pm, you would create two separate files, maybe outside_lights.yaml and sprinklers.yaml

###File Type (yaml)
All automation must be in a file that has a ".yaml" file extension. If the file does not have this extension goHOME will not attempt to load it. For example, your automation file could be called sunset.yaml which contains all of the automation you want to execute when it is sunset.

###What is yaml?
yaml is a compact and human friendly way to describe data.  If you've never written it before, the parts we use in goHOME are very simple, take a quick look at this page: https://learnxinyminutes.com/docs/yaml/ for our purposes, you just need to be able to write comments, undestand keys and values and lists.

###Finding errors in your script
When writing your automation, you may have errors in your script.  To check if your script is valid, save the file, then restart the gohome executable, if you look at the app log on startup it will say something like:
"automation - [script name]" 
if it loads successfully. It is fails to load there will be an error written to the output.

###Testing Automation
When you are writing some automation, rather than having to wait until the trigger fires to test your script to make sure it executes as expected, you can test the automation and make it execute immediately.  Once you have written the file, restart the goHOME server and the new script will be loaded, no in the UI, click on the "automation" tab in the app header, you will see your automation listed in the UI. IF you click on the item, a "Test" button will appear, clicking on it will immediately execute your automation, so you can verify it is working as expected.

![](docs/automation.png)

###Syntax
Here is an example automation script, lets call it sunset.yaml More details on the exact syntax and all allowable values are listed after this example.

```yaml
name: 'Sunset'
trigger:
  time:
    at: sunset
    days: sun|mon|sat
actions:
  - window_treatment:
      open_closed: 'closed'
  - light_zone:
      # Entry
      id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
      on_off: 'on'
      brightness: 100
  - light_zone:
      # Front Door
      id: ea92dae8-cfba-4fc3-57df-8f7e13d231fc
      on_off: 'on'

```

Lets go through this piece by piece. IMPORTANT: to get the IDs for the features you wish to control, go to the "feature" tab in the app and click on the edit icon in the top right, the view will change and you will see the ID listed under each feature.

```yaml
name: 'Sunset'
```
This specifies the name of the automation. In yaml you can specify key:value pairs, strings are enclosed inbetween single quotes '' (they don't always have to be, but to save some pain, just always use them).  All automation scripts need a name value.

```yaml
trigger:
  time:
    at: sunset
    days: sun|mon|sat
```
The next part defines the automation trigger. There can be different types of triggers, e.g. time, sensor change, button pressed etc. For this example we have a time trigger that can be used to trigger events at certain times. IMPORTANT: notice the indentation, it is very important, you must use spaces (not tabs) to indent the different pieces.

In this case our time trigger will fire at sunset (this time varies from day to day and is defined by the location setting in your config file).  We also have an optional "days" field which means this trigger will only fire on Sunday, Monday and Saturday. If you don't include the "days" key then the trigger will fire every day.

```yaml
actions:
  - window_treatment:
      open_closed: 'closed'
  - light_zone:
      # Entry
      id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
      on_off: 'on'
      brightness: 100
  - light_zone:
      # Front Door
      id: ea92dae8-cfba-4fc3-57df-8f7e13d231fc
      on_off: 'on'
```

Next we define the actions we wish to execute when the trigger fires. Note the '-' character, in yaml that denotes a list item, you can have one or more actions execute. In this example you see we execute three actions. These actions will be executed sequentially.  Each action has a type of feature it affects, in the example we see window_treatment and light_zone.

For the window_treatment we want to close all of the window shades in our house at sunset, so note we don't specify a specific ID, which means this action applies to all window treatments in the system. We set the open_closed state to 'closed'.

For the light_zone actions, we want to control specific lights, in this case at sunset I want the lights outside the front door and the entry light to come on, so it is bright when we come home from work. Note we set the on_off state and the brightness state.  For the front door light, it is an on/off light it can't be dimmed so we don't specify brightness.  Like the window treatment, you could omit the ID field and then the action would be applied to all light zones.

##More Examples
For more examples see the following repositories: //TODO:

##Detailed Syntax
Here we list the complete automation syntax.

##Triggers
The following triggers are supported:

###Time Trigger
The create a time trigger, add the "time" key to the trigger:
```yaml
trigger:
  time:
```
The following fields are supported on the time trigger:
####at (required)
Values: sunset|sunrise|yyyy/MM/dd HH:mm:ss|HH:mm:ss

  - sunset -> The time trigger will fire at sunset (as defined by the Location value in your config.json file)
  - sunrise -> The time trigger will fire at sunrise
  yyyy/MM/dd HH:mm:ss -> specifies an exact date and time the trigger should fire. The trigger will only fire once on this exact datetime, the time needs to be in 24 hour format and always include the seconds e.g. 2016/10/28 19:40:00
  - HH:mm:ss -> specifies a time for the trigger to execute. Note we don't specify the date, so the trigger will fire every day at this time (see "days" field for more info on how to change this)

####days (optional)
Values: sun|mon|tues|wed|thurs|fri|sat

If you don't specify a "days" key then the trigger fires every day (as long at the time was not specified with a date and time). You can specify any number of days separated by a | character. For example, to specify the trigger should fire on Tuesday and Friday you would use the value tues|fri
 
##Actions
There are many actions we can execute when a trigger is fired, below are the complete list
###light_zone
Turns lights on/off or to specific brightnesses (if supported)
```yaml
light_zone:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
  on_off: 'on'
  brightness: 39
```
####id (optional)
The id of the specific light zone you wish to control. If you exclude this key, the action will be applied to "all" light zones in the system.
####on_off (optional)
Values: 'on'|'off'
IMPORTANT: Make sure you include the single quotes around the values, otherwise your script will not load.
####brightness (optional)
A value between 0 and 100. If you specify this value on a light that doesn't support dimming it will be ignored.

###switch
Turns a switch on or off
```yaml
switch:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
  on_off: 'off'
```
####id (optional)
Specifies the id of the switch to control, if omitted this action is applied to all switches.
####on_off (required)
Values: 'on'|'off'

###outlet
Turns an outlet on or off
```yaml
outlet:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
  on_off: 'on'
```
####id (optional)
Specifies the id of the outlet to control, if omitted this action is applied to all outlets.
####on_off (required)
Values: 'on'|'off'

###window_treatment
Controls the offset of a window treatment such as a shade
```yaml
window_treatment:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
  offset: 75
  open_closed: 'open'
```
####id (optional)
The id of the window treatment to control, if not specified the action is applied to all window treatments
####open_closed
Values: 'open'|'closed'
Specifies if the window treatment is open or closed. If closed the offset parameter is ignored, if open and no offset parameter is specified the window treatment will open to 100%
####offset (optional)
A value between 0 and 100. 0 represents fully closed and 100 is fully open.

###heat_zone
Controls a heat zone temperature
```yaml
heat_zone:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
  target_temp: 75
```
####id (optional)
The id of the heat zone to control, if ommitted the action is applied to all heat zones.
####target_temp (required)
A value between 40 and 80, representing the target temperature to set in Farenheit.

###scene
The scene action executes the specified scene.
```yaml
scene:
  id: 3fc087bc-0660-4aec-7a55-f5259a5b4119
```
####id (required)
The id of the scene to execute
