
# WMTUI

[comment]: <> (General About Section)
### About:

**WMTUI** was designed as an alternative to WhatsminerTool to be used in the terminal of the local master server of mining sites.

[comment]: <> (Basic Features)
### Features

- See Whatsminers on a set ip range
  - view: IP, Mac, Status, Error Codes, Uptime, GHs, Efficiency, Wattage, Set Power Limit, 1st pool URL
- Switch between multiple saved ip ranges
- Send command to one or many machines
    - Sleep and Wake machines
    - Set pools
    - Set Power Limit
    -  Enable / Disable Fastboot
    -  Reboot

[comment]: <> (Keycommands Laid out)
### Keyboard Commands

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

[comment]: <> (How to Generaly use WMTUI)
### Basic Usage

- Add a farm by with `n`, filling out the details and pressing `enter`.
    - After creating a new farm you enter the farm selection screen where you can select your newly created farm with `enter`.
- Farms can be switched between with `a`, and selected with `enter`.
- Press `m` to refresh the list of miners and miner info.
- Miners can be selected with `enter` which will remove them from the main page and send them to the selected page.
    - Miners can also be removed from the selected screen with `enter`.
- Selected miners can be seen by pressing `tab` all selected miners will receive command while on this screen.
- Commands sent from the main page are only sent to the highlighted machine.

[comment]: <> (List of bugs and missing features)
### Known Bugs and Missing Features

- The view can break if your terminal window is too small.
- The view will break if resizing the terminal window.
- Setting pools has no validation check on both the pool URL and worker name.
- Setting farms has no IP validation.
- Setting powerlimit has no validation.
- Farm IP ranges are only based on the last 2 octets of the IP address.
- API Errors are currently unhandled and may crash the program.
- No feedback on command success or failure.
- Scanning for miners will occasionally miss some.
- Selected miners can be sent back to the main page as duplicates if refreshing after selecting.
  - Refresh again to get rid of them.
- There's currently no way to bulk remove selected machines
- WMTUI will occasionally lock up on launch requiring a force killing of the main process.
- WMTUI will occasionally crash on startup.
- When switching between farms the rescan doesn't always scan properly, requiring a restart to fully switch to the new farm.

[comment]: <> (How to Install From Release)
### Installing from release

Download the appropriate release version from the release page and unzip it. </br>
Open a terminal and navigate to the newly created folder and run one of the following commands:</br>
> Windows: </br>
> `.\WMTUI-[arcitecture].exe`

> Mac/Linux: </br>
> `./WMTUI-[arcitecture]`

[comment]: <> (How to build for src)
### Building From Source

#### Requires go >= v1.22</br>

Download the file from the `<> code` button in zip format and unzip it. </br> 
From the root directory of the project run the following: </br>
`go mod tidy` </br>
`go build .`

[comment]: <> (How to add a farm outside of WMTUI)
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

