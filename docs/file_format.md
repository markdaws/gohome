Your system configuration is stored as JSON. Below is a sample configuration file, you can modify your config file by hand if the UI does not support the action you are trying to perform, but always make a backup of your last known good file.

```json5
{
  // The version of the configuration file format
  "version": "0.1.0",
  
  "name": "Lutron Smart Bridge Pro",
  "description": "Lutron Smart Bridge Pro",
  
  // Used internally to generate global IDs - do NOT edit
  "nextGlobalId": 210,
  
  // Array of scene objects
  "scenes": [
    {
      "address": "2",
      "id": "103",
      "name": "All Off",
      "description": "All Off",
      "commands": [
        {
          "type": "buttonPress",
          "attributes": {
            "ButtonID": "52"
          }
        },
        {
          "type": "buttonRelease",
          "attributes": {
            "ButtonID": "52"
          }
        }
      ]
    }
  ],
  "devices": [
    {
      "address": "192.168.0.10:23",
      "id": "114",
      "name": "Smart Bridge - Hub",
      "description": "",
      "modelNumber": "l-bdgpro2-wh",
      "modelName": "",
      "softwareVersion": "",
      "hubId": "",
      "buttons": [
        {
          "address": "41",
          "id": "41",
          "name": "Button 41",
          "description": "Button 41"
        }
      ],
      "zones": [
        {
          "address": "23",
          "id": "198",
          "name": "Living Room Shade",
          "description": "Living Room Shade",
          "deviceId": "114",
          "type": "shade",
          "output": "continuous"
        },
        {
          "address": "16",
          "id": "194",
          "name": "Dining Area",
          "description": "Dining Area",
          "deviceId": "114",
          "type": "light",
          "output": "continuous"
        },
      ],
      "sensors": [
        
      ],
      "deviceIds": [
        "144",
        "138",
      ],
      "auth": {
        "login": "bob",
        "password": "12345",
        "token": ""
      },
      "connPool": {
        "name": "Smart Bridge - Hub",
        "poolSize": 2
      }
    },
    {
      "address": "http:\/\/192.168.0.34:49154",
      "id": "204",
      "name": "WeMo Maker abcd",
      "description": "Belkin Plugin Socket 1.0",
      "modelNumber": "1.0",
      "modelName": "Maker",
      "softwareVersion": "WeMo_WW_2.00.6529.PVT",
      "hubId": "",
      "buttons": [
        
      ],
      "zones": [
        {
          "address": "1",
          "id": "205",
          "name": "WeMo Maker abcd",
          "description": "Belkin Plugin Socket 1.0",
          "deviceId": "204",
          "type": "switch",
          "output": "binary"
        }
      ],
      "sensors": [
        {
          "id": "206",
          "name": "WeMo Maker abcd - sensor",
          "description": "",
          "address": "1",
          "deviceId": "204",
          "attr": {
            "name": "sensor",
            "value": "",
            "dataType": "int",
            "states": {
              "0": "Closed",
              "1": "Open"
            }
          }
        }
      ],
      "deviceIds": [
        
      ],
      "auth": null,
      "connPool": null
    }
  ]
}
```
