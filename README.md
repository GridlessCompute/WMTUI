
# WMTUI

### About:

**WMTUI** was designed as an alternative to WhatsminerTool to be used in the terminal of the local master server of mining sites.

### Features:

- See Whatsminers on a set ip range
  - view: IP, Mac, Status, Error Codes, Uptime, GHs, Efficiency, Wattage, Set Power Limit, 1st pool URL
- Switch between multiple saved ip ranges
- Send command to one or many machines
    - Sleep and Wake machines
    - Set pools
    - Set Power Limit
    -  Enable / Disable Fastboot
    -  Reboot

### Keyboard Commands:

- `↑` &nbsp;&nbsp; - Move up the miner list
- `↓` &nbsp;&nbsp; - Move down the miner list
- `'` &nbsp;&nbsp; - Sort miners, toggles between four modes, ip, mac, hash rate, and uptime. Ip sort is the default
- `?` &nbsp;&nbsp; - Show help menu
- `tab`&nbsp;&nbsp; - Switch Between the miner list pane and the selected miners pane
- `enter`&nbsp;&nbsp; - Select or deselect a miner
- `a`&nbsp;&nbsp; - Enter the farm selection screen
- `f`&nbsp;&nbsp; - Enable fastboot 
- `l`&nbsp;&nbsp; - Set power limit
- `m`&nbsp;&nbsp; - Refresh the miner view
- `n`&nbsp;&nbsp; - Create a new farm
- `o`&nbsp;&nbsp; - Disable fastboot 
- `p`&nbsp;&nbsp; - Set pools
- `q`&nbsp;&nbsp; - Quit 
- `r`&nbsp;&nbsp; - Reboot 
- `s`&nbsp;&nbsp; - Sleep 
- `w`&nbsp;&nbsp; - Wake

### Known Bugs and Missing Features

- The view can break if your terminal window is too small
- Setting pools has no validation check on both the pool URL and worker name.
- Setting farms has no IP validation
- Farm IP ranges are only based on the last 2 octets of the IP address
- API Errors are currently unhandled and may crash the program.
- No feedback on command success or failure.
- Scanning for miners will occasionally miss some.
- WMTUI will occasionally lock up on launch requiring a force killing of the main process.
- WMTUI will occasionally crash on startup
- When switching between farms the rescan doesn't always scan properly, requiring a restart to fully switch to the new farm.

### Manual Farm Setup
 WMTUI comes preset with the `192.168.0.0/24` ip range and more can be added while running.  
 The `FARMS.JSON` file can be manually edited as well if preferred,  
 farms are stored in the following JSON format:
 >{<br>
 > &nbsp;&nbsp;&nbsp;&nbsp;"Farms": [{<br>
 >	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"Name": "Name of farm"<br>
 >	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"Start": "192.168.0.1"<br>
 >	&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"End": "192.168.0.255"<br>
 > &nbsp;&nbsp;&nbsp;&nbsp;}]<br>
 >}

