During development of this project, there are many hardware quirks that have been discovered.  This document will list all of the issues that we found.

##FluxWIFI Bulbs
  - Getting current state is not always accurate.  After setting a new RGB value, querying the bubl for the current values still returns old values, at some point these values update but it can take a while.
  
##Lutron Smart Bridge Pro
  - The bridge can stop responding for 30 seconds to a minute after receiving many commands.  Try to reduce the number of commands sent to this device in a small period of time.  There seems to be some kind of rate limiting in place internally.
