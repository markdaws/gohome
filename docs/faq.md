##How do I...

###Change a light zone to only turn on/off and not be dimmable?
Ideally when you import light zones we can tell if they are on/off only or dimmable, but some systems don't give you enough information. You can't change this via the UI at the moment, to remove the dimmable functionality:
  - Go to the features tab, and click on the edit button in the top right
  - Find the light zone you wish to change, note down the ID of the light zone
  - In a text editor open your gohome.json file (or whatever you called your configuration file)
  - Find the section keyed by "features: [", scroll down to the feature that matches your ID
  - Find the "attrs" key and remove the "brightness" key located in that
  - Restart your system, the light will now no longer be dimmable
