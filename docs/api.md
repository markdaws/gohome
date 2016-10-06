The following HTTP API is exposed by the goHOME server:

##Devices
//TODO:

##Scenes
###GET /api/v1/scenes
Returns an array containing all of the scenes in the system

#### Example Response
 ```json
 {
    "address": "2",
    "id": "215",
    "name": "All Off",
    "description": "Turns all the lights off",
    "managed": false,
    "commands": [
      {
        "type": "buttonPress",
        "attributes": {
          "ButtonID": "53"
        }
      },
      {
        "type": "buttonRelease",
        "attributes": {
          "ButtonID": "53"
        }
      }
    ]
  },
 ]
 ```
 
 
 - GET /api/v1/scenes/&lt;ID&gt;
 - DELETE /api/v1/scenes/&lt;ID&gt;
 - PUT /api/v1/scenes/&lt;ID&gt;
 - POST /api/v1/scenes
 - POST /api/v1/scenes/active
 - POST /api/v1/scenes/&lt;ID&gt;/commands
 - DELETE /api/v1/scenes/&lt;ID&gt;/commands/&lt;INDEX&gt;

##Zones
 - GET /api/v1/zones
 - POST /api/v1/zones
 - PUT /api/v1/zones/&lt;ID&gt;

##Discovery
 - GET /api/v1/discovery/&lt;MODEL_NUMBER&gt;
 - GET /api/v1/discovery/&lt;MODEL_NUMBER&gt;/token
 - GET /api/v1/discovery/&lt;MODEL_NUMBER&gt;/access
 - GET /api/v1/discovery/&lt;MODEL_NUMBER&gt;/zones
