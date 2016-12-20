<p align="center">
<img src="logo.png" />
</p>
goHOME is an open source home automation client + server, aimed at allowing home owners to have full control over their home automation hardware. It provides a single fully open and customizable UI to control many different pieces of home automation hardware under one UI. The project also runs on cheap hardware like a Raspberry PI.

Note - still in alpha development, do not use just yet :)

//TODO: Pics

##Why use an open source home automation project?
###Unified User Interface
I have home automation devices from many different manufacturers in my home, each device requires a different control app, either logging in to a browser app, or installing a mobile phone app, so to control my home I have to use serveral apps. Having one consistent UI where you can control all your devices is a big win for ease of use.  Also, software is not a primary focus of hardware manufacturers, the UIs to control these devices is generally pretty terrible, an after thought. With goHOME the aim is to provide a great user experience and allow the community to customize and provide enahanced experiences not available in the original manufacturer apps.

###Local control - no outside services
One of the main aims of the project is to interop with the hardware on your local network.  Too many home automation devices now require you to connect to an outside service just to talk to the device on your local network.  This is inefficient and error prone.  I have had serveral times where I wanted to change my lights but couldn't because the external service was down, even though the hub on my local network was fine.

###Interoperability
You want all of your devices to be able to work with one another.  By providing an additional unification layer on top of all the hardware devices it is easy for us to plug devices together, without requiring support from the manufacturers.

###Device/Service obsolescence
If you have invested money in buying these hardware devices, it sucks when a manufacturer decides to shut down the service that controls the device or no longer provide support for the control app etc.  By using goHOME the aim is that we can make sure you can control and monitor your devices forever, even long after the manufacturer has given up providing support. 

##Supported Hardware
For a list of the supported hardware, click <a href="docs/supported_hardware.md">here</a>.  If you want support for a piece of hardware that is not on the list, file an issue.

##Documentation
//TODO: Core
//TODO: Architecture
###[Getting Started](docs/getting_started.md)
###[Supported Hardware](docs/supported_hardware.md)
###[Raspberry PI Setup](docs/raspberry_pi.md)
###[FAQ](docs/faq.md)
###[Automation](docs/automation.md)

##Development
Currently there are two main parts to the project, the golang powered server and the web UI.

###goHOME Server
The goHOME server is developed using golang (http://golang.org) In order to develop the code:

  - Install git (source control): https://git-scm.com/
  - Install golang https://golang.org/dl/
  - Setup your GOPATH https://golang.org/doc/code.html#GOPATH

Once you have done this, run the following commands:
```bash
go get github.com/markdaws/gohome
```

Change to the source directory which will be $GOPATH/src/github.com/markdaws/gohome and build the app, running:
```bash
go build -o gohome ./main
```

####Adding a user account
You need to add a user to be able to log into the app, for example we will add a user "bob" with password "foobar" by running the gohome executable in the source directory $GOPATH/src/github.com/markdaws/gohome

```bash
./gohome --set-password bob foobar
```

####Starting the server
Then start the server:
```bash
./gohome --server
```
In the output you will see a line like (note the IP address is probably different):
```
WWW Server starting, listening on 192.168.0.10:8000
```

###gohome web UI
The web UI is developed using the React framework: https://facebook.github.io/react/ In order to develop the web UI:
 1. Setup the goHome Server, following the above instructions
 2. Install node.js: https://nodejs.org
 3. Change to the gohome/www directory
 4. Run:
 
 ```bash
 npm install
 ```
 5. Run webpack to monitor file changes and build the UI:
 
 ```bash
 npm run dev
 ```

####NOTE - All web UI code is located at gohome/www/assets
####NOTE - IF you want a production build of the UI, with minified source, run
```bash
npm run prod
```
