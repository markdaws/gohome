Currently there are two main parts to the project, the golang powered server and the web UI.  If you only want to change the server code, it comes with all the necessary UI files in the build, so you don't need to build the UI, just build and run the server.

##goHOME Server
The goHOME server is developed using golang (http://golang.org) In order to develop the code:

  - Install git (source control): https://git-scm.com/
  - Install golang https://golang.org/dl/
  - Setup your GOPATH https://golang.org/doc/code.html#GOPATH

Once you have done this, run the following commands:
```bash
go get github.com/markdaws/gohome
```
You may see some output like the following, you can ignore it:
```bash
package github.com/markdaws/gohome: no buildable Go source files in /home/gohome/go/src/github.com/markdaws/gohome
```

Change to the source directory which will be $GOPATH/src/github.com/markdaws/gohome and build the goHOME executables, running:
```bash
./build.sh
```
Once the build has completed, you will see two binaries in the $GOPATH/bin directory:
  - ghadmin: An admin tool for creating a new project and adding new users
  - ghserver: The goHOME server executable
  
###Creating your project
goHOME needs tow main files, a config file that specifies settings, like which IP address and port to use for the web server and a system file that contains all of your pject information, such as which hardware you have imported, user information etc. The first thing we have to do is init these files, choose a directory somewhere that you want to store these files, then run the following command (the argument specifies the directory where the goHOME source code is located):
```bash
ghadmin --init $GOPATH/src/github.com/markdaws/gohome
```
After the command runs, in the current directory you will see a config.json and gohome.json file, take a look inside. If there are any settings you want to change in config.json you can make them now.

###Adding a user account
You need to add a user to be able to log into the app, for example we will add a user "bob" with password "foobar", you have to specify the location of the config.json file that was created in the previous step:

```bash
ghadmin --config=/path/to/my/config.json --set-password bob foobar
```

###Starting the server
The server is responsible for communicating with all of your home automation hardware and serving the web UI. To start the server:
```bash
ghserver --config=/path/to/my/config.json
```
In the output you will see a line like (note the IP address is probably different):
```
WWW Server starting, listening on 192.168.0.10:8000
```

##goHOME web UI
The web UI is developed using the React framework: https://facebook.github.io/react/ In order to develop the web UI:
 1. Setup the goHOME Server, following the above instructions
 2. Install node.js: https://nodejs.org
 3. Change to the gohome/pkg/www directory
 4. Run:
 
 ```bash
 npm install
 ```
 5. Run webpack to monitor file changes and build the UI:
 
 ```bash
 npm run dev
 ```

####NOTE - All web UI code is located at gohome/pkg/www/assets
####NOTE - If you want a production build of the UI, with minified source, run
```bash
npm run prod
```

####NOTE - when you build the web UI, all the files are copies and served from the root "dist" folder, not the gohome/pkg/www folder.
