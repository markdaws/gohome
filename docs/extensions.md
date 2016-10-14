New device support can be added to goHOME via an extension architecture.  An extension provides all of the device specific know how which is needed to control and monitor a device.  All extensions are located under the extensions folder.

##Adding a new extension
If you want to add a new extension, awesome! Please follow the guidelines below:

###Submit an issue before spending time updating the code
Before you fork the code and submit a pull request, it's always best to submit an issue indicating what hardware you are going to support and roughly how you will modify the code.  This can save everyone time, incase the changes you plan to make don't match how the code should be modified.

###Code Modifications
  1. Add a new folder under the extensions folder with a name that represents the class of devices the extension controls e.g. "lutron", "honeywell"
  
  2. Add a new file called extension.go to the folder, this is the entry point where the functionality is registered
  
  3. Update the RegisterExtensions function inside /master/intg/intg.go to register your new extension.  The code is pretty self explanatory.
  
  //TODO: inteface examples etc

//TODO: run gofmt

###Testing
//TODO: Add unit test requirements

###Documentation
//TODO: List the docs that should be updated

###Submitting
Once you have updated your forked version of the code to include the new extension, submit the code as a pull request and we will take a look.
