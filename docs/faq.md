### How do I change a light zone to only turn on/off and not be dimmable?
Ideally when you import light zones we can tell if they are on/off only or dimmable, but some systems don't give you enough information. You can't change this via the UI at the moment, to remove the dimmable functionality:
  - Go to the features tab, and click on the edit button in the top right
  - Find the light zone you wish to change, note down the ID of the light zone
  - In a text editor open your gohome.json file (or whatever you called your configuration file)
  - Find the section keyed by "features: [", scroll down to the feature that matches your ID
  - Find the "attrs" key and remove the "brightness" key located in that
  - Restart your system, the light will now no longer be dimmable

### Why are my sunrise/sunset events at the wrong time
You need to specify the latitude/longitude of the gohome server in the [config](config.md). Also if you are using a Raspberry PI make sure you have the correct time zone set on your device, to change it run the following command:
```bash
sudo dpkg-reconfigure tzdata
```

### Why can't I log in after installing goHOME
goHOME doesn't have a default user account, you need to create one, or maybe you forgot your password, see [getting started](getting_started.md) for more details
