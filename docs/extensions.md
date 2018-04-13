//TODO:
Work in progress...

New device support can be added to goHOME via an extension architecture.  An extension provides all of the device specific know how which is needed to control and monitor a device.  All extensions are located under the extensions folder.

## Adding a new extension
If you want to add a new extension, awesome! Please follow the guidelines below:

### Submit an issue before spending time updating the code
Before you fork the code and submit a pull request, it's always best to submit an issue indicating what hardware you are going to support and roughly how you will modify the code.  This can save everyone time, incase the changes you plan to make don't match how the code should be modified.

### Extension Interface
All extensions need to satisfy the following interface, in order to be considered as a valid extensions.  We will go through each one of these methods in more details below:

```go
type Extension interface {
	Name() string
	BuilderForDevice(*System, *Device) cmd.Builder
	NetworkForDevice(*System, *Device) Network
	EventsForDevice(*System, *Device) *ExtEvents
	Discovery(*System) Discovery
}
```

#### Name()
Should return a unique name for your extension - this name will be printed in the system log, in debug information.

####BuilderForDevice(sys *System, dev *Device) cmd.Builder
Since goHOME abstracts on top of many different types of hardware, there are abstract commands that the app uses, such as ZoneSetLevel, SetScene, ButtonPress etc, that know nothing about how to actually control specific hardware.  For this function you are passed in a device, you should check to see if you own this device and know how to build commands for it, if not, just return nil.  If you own the device then you just need to return a type that satisfies the cmd.Builder interface, which looks like:
```go
type Builder interface {
	Build(Command) (*Func, error)
}
```

### NetworkForDevice(sys *System, dev *Device) Network
//TODO:
### EventsForDevice(sys *System, dev *Device) ExtEvents
//TODO:
### Discovery(sys *System) Discovery
//TODO:

### Example Extension
There is a basic example extension under the gohome/extensions/example folder, you can copy this extension into your new folder and update it for your specific device.

### Best practices
  - Always set timeouts for network read/write operations, using SetReadDeadline/SetWriteDeadline
  - Write defensive code, don't expect bad things not to happen
  - Always check error values and react appropriately
  - Gracefully recover from failures. If a network connection closes, try to open a new one, don't sit in a tight loop though, add some sleep periods between attempting to connect.
  
### Code Modifications
  1. Add a new folder under the extensions folder with a name that represents the class of devices the extension controls e.g. "lutron", "honeywell"
  
  2. Add a new file called extension.go to the folder, this is the entry point where the functionality is registered
  
  3. Update the RegisterExtensions function inside /master/intg/intg.go to register your new extension.  The code is pretty self explanatory.
  
  //TODO: inteface examples etc

//TODO: run gofmt

### Testing
//TODO: Add unit test requirements

### Documentation
//TODO: List the docs that should be updated

### Submitting
Once you have updated your forked version of the code to include the new extension, submit the code as a pull request and we will take a look.
