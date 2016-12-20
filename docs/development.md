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

Change to the source directory which will be $GOPATH/src/github.com/markdaws/gohome and build the app, running:
```bash
go build -o gohome ./main
```

###Adding a user account
You need to add a user to be able to log into the app, for example we will add a user "bob" with password "foobar" by running the gohome executable in the source directory $GOPATH/src/github.com/markdaws/gohome

```bash
./gohome --set-password bob foobar
```

###Starting the server
Then start the server:
```bash
./gohome --server
```
In the output you will see a line like (note the IP address is probably different):
```
WWW Server starting, listening on 192.168.0.10:8000
```

##goHOME web UI
The web UI is developed using the React framework: https://facebook.github.io/react/ In order to develop the web UI:
 1. Setup the goHOME Server, following the above instructions
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
####NOTE - If you want a production build of the UI, with minified source, run
```bash
npm run prod
```
