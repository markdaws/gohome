Below are all of the config settings for the goHOME application

```json5
{
  //The full path to the file that contains all of your system configuration, such as lights, shades etc
  //By default if not set gohome creates a file called gohome.json in the same directory as the gohome executable
  systemPath: "",

  //The full path to where the event log will be written. By default a file called events.json is create in the 
  //same directory as the gohome executable
  eventLogPath: "",

  //The path where goHOME will look for your automation scripts. By default it will look for a directory called
  //"automation" in the directory where the gohome executable is located
  automationPath: "",

  //The IP address for the WWW server. By default gohome looks for the first non loopback address
  wwwAddr: "",

  //The port to use for the WWW server, defaults to "8000"
  wwwPort: "",

  //The IP address for the API server. By default gohome looks for the first non loopback address
  apiAddr: "",

  //The port to use for the API server, defaults to "5000"
  apiPort: "",

  //The IP address used for a UPNP notify server, gohome looks for the first non loopback address
  upnpNotifyAddr: "",

  //The port used for the UPNP notification server, defaults to "5001"
  upnpNotifyPort: "",

  //If you want sunset/sunrise events to have the correct time, you have to specify the location where the 
  //gohome server is located
  location: {
    latitude: 0.0,
    longitude: 0.0
  }
}
```
