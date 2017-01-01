The Raspberry PI is an excellent candidate for an inexpensive home automation server.  For an application like goHOME you want a computer that is running 24/7 connected to your network, if you have a desktop computer that you always leave on then you can install goHOME on that, but if you want a dedicated device, for around $25 you can buy a Raspberry PI.

##Setting up your Raspberry PI
There are many excellent tutorials on how to install an OS on your Raspberry PI, here is one that shows how to easily install Raspbian (the preferred Raspberry PI Operating System): https://www.raspberrypi.org/learning/software-guide/quickstart/ Once your Raspberry PI is up and running, you can move on to the next stage.

##Connecting to your Raspberry PI (via ssh)
If you have a keyboard and mouse connected to your Raspberry PI then you can skip this step, otherwise we will want to use ssh to connect to the Raspberry PI and set it up.

IMPORTANT: In order to connect via ssh you will need to enable SSH on your Raspberry PI. To do this create a file called "ssh" on the root partition on your SD card where you installed Raspbian, if you don't do this you will get a connection refused error when you try to ssh.

###Finding the IP address of your Raspberry PI
We need the IP address of your Raspberry PI in order to connect to it, see this tutorial https://www.raspberrypi.org/documentation/remote-access/ip-address.md

For all of the following, I am using the IP address 192.168.0.37, but you should replace that with whatever IP address your device was assigned.

###Connecting via SSH
```bash
ssh pi@192.168.0.37
```
The default password is "raspberry"

Once connected change the default password of the pi user, don't leave it as the default
```bash
passwd
```

##Updating your OS
Make sure you have the latest software on your Raspberry PI
```bash 
sudo apt-get update
sudo apt-get upgrade -y
```

##Creating a goHOME user
We will create a dedicated user for the goHOME installation.
```bash
sudo useradd -rm gohome
```
Then set the password for the gohome user
```bash
sudo passwd gohome
```

##Installing git
git is used to download the latest goHOME source code, to install run:
```bash
sudo apt-get install git
```

##Installing golang
[golang](https://golang.org/) is the language used to write the goHOME server. You need to install the go runtime to run goHOME on your Raspberry PI. Follow the instructions below:

```bash
wget https://storage.googleapis.com/golang/go1.7.4.linux-armv6l.tar.gz
sudo tar -xzf go1.7.4.linux-armv6l.tar.gz -C /usr/local
sudo chgrp -R staff /usr/local/go
export GOROOT=/usr/local/go
export PATH="$PATH:$GOROOT/bin"
```
NOTE: There is a version of go available via apt, but don't use it because it is an old version and goHOME won't run, use the instructions above.

At this point you can test your go installation by running the following and making sure you get a response:
```bash
go version
```

##Initializing the environment
Switch to the gohome user
```bash
su gohome
```
We need to make a directory where the goHOME source will be located, create one called "go" in the /home/gohome user directory:
```bash
cd ~
mkdir go
```

Next we need to define some environment variables so that the go runtime can build and execute our program, instead of doing this every time we log in, we will add them to /home/gohome/.bashrc so they are setup each time you log in
```bash
nano ~/.bashrc
```
At the end of the file add the env variables:
```bash
export GOROOT=/usr/local/go
export PATH="$PATH:$GOROOT/bin"
export GOPATH=/home/gohome/go
export PATH="$PATH:$GOPATH/bin"
```
At this point we will log out of the SSH session and log back in as the gohome user
```bash
ssh gohome@192.168.0.37
```

##Download the goHOME source code
We will download the source code to the goHOME application
```bash
go get github.com/markdaws/gohome
```

##Building the application
Change to the goHOME source directory
```bash
cd /home/gohome/go/src/github.com/markdaws/gohome
```
Then build the applications
```bash
./build.sh
```

Once the build has completed, you will see two binaries in the $GOPATH/bin directory:
  - ghadmin: An admin tool for creating a new project and adding new users
  - ghserver: The goHOME server executable
  
##Creating your project
goHOME needs tow main files, a config file that specifies settings, like which IP address and port to use for the web server and a system file that contains all of your pject information, such as which hardware you have imported, user information etc. The first thing we have to do is init these files, choose a directory somewhere that you want to store these files, then run the following command (the argument specifies the directory where the goHOME source code is located):
```bash
ghadmin --init $GOPATH/src/github.com/markdaws/gohome
```
After the command runs, in the current directory you will see a config.json and gohome.json file, take a look inside. If there are any settings you want to change in config.json you can make them now.

##Adding a user account
You need to add a user to be able to log into the app, for example we will add a user "bob" with password "foobar", you have to specify the location of the config.json file that was created in the previous step:

```bash
ghadmin --config=/path/to/my/config.json --set-password bob foobar
```

NOTE: If you are putting special chars in your password, enclose it in double quotes so that it is not interpretted by the shell.

##Running the app
Finally, phew. Now we can launch the gohome app, run
```bash
ghserver --config=/path/to/my/config.json &
```

You will see the app output some information to the terminal, you should see a line like:
```bash
WWW Server starting, listening on 192.168.0.37:8000
```

You can now go to your browser and load the app at the specified location.

##Changing the IP address/port 
If you don't want the default IP/port addresses, you can change the values in the [config](docs/config.md)

##Running the app on boot
Before running the following commands, change back to the "pi" user instead of the gohome user.

To automatically run the gohome server when the Raspberry PI reboots, create a file at the following location /etc/systemd/system/gohome.service and paste the following contents

```
[Unit]
Description=goHOME Server
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
ExecStart=/home/gohome/go/ghserver --config=/home/gohome/go/src/github.com/markdaws/gohome/config.json
Restart=always
User=gohome
WorkingDirectory=/home/gohome/go/src/github.com/markdaws/gohome

[Install]
WantedBy=multi-user.target
```
Save the file, then run the following to make the goHOME server run on boot
```bash
sudo systemctl enable gohome
```

You can start the service by running, without rebooting by
```bash
sudo systemctl start gohome
```
