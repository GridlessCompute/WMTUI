



<p align="center">
  <img src="https://github.com/GridlessCompute/WMTUI/assets/98416347/3c8eb6ea-8d13-4caf-bf1b-7820b8aae790" alt="Sublime's custom image"/>
</p>
<p align="center">
  <img src="https://github.com/GridlessCompute/WMTUI/assets/98416347/6c61100f-22df-42e9-b9c2-9e12b72653cb" alt="Sublime's custom image"/>
</p>


# WMTUI

[comment]: <> (General About Section)
### About:

**WMTUI** was designed as an alternative to WhatsminerTool to be used on the terminal of the local server of mining sites.

[comment]: <> (Basic Features)
### Features

- View Whatsminers on a set ip range
  - view: IP, Mac, Status, Error Codes, Uptime, GHs, Efficiency, Wattage, Set Power Limit, 1st pool URL
- Switch between multiple saved ip ranges
- Send command to one or many machines
    - Sleep and Wake machines
    - Set pools
    - Set Power Limit
    - Enable / Disable Fastboot
    - Reboot

[comment]: <> (Keycommands Laid out)
### Keyboard Commands

#### Miner List / Selected Miner List

- `esc` &nbsp;&nbsp; - Quit
- `↑` &nbsp;&nbsp; - Move up the miner list
- `↓` &nbsp;&nbsp; - Move down the miner list
- `'` &nbsp;&nbsp; - Sort miners, toggles between four modes, ip, mac, hash rate, and uptime. Ip sort is the default
- `?` &nbsp;&nbsp; - Show help menu
- `tab`&nbsp;&nbsp; - Switch between the miner list pane and the selected miners pane
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

#### New Farm

- `esc` &nbsp;&nbsp; - Go back to the miner list
- `tab`&nbsp;&nbsp; - Switch Between text fields
- `enter`&nbsp;&nbsp; - Save new farm

#### Select Farm

- `esc` &nbsp;&nbsp; - Go back to the miner list
- `↑` &nbsp;&nbsp; - Move up the farm list
- `↓` &nbsp;&nbsp; - Move down the farm list
- `enter`&nbsp;&nbsp; - Select farm

#### Set Pools

- `esc` &nbsp;&nbsp; - Go back to the miner list
- `tab`&nbsp;&nbsp; - Switch Between text fields
- `enter`&nbsp;&nbsp; - Send the new pools

#### Set Power Limit

- `enter`&nbsp;&nbsp; - Send the new power limit

[comment]: <> (How to Generaly use WMTUI)
### Basic Usage

- Add farm - Add a new farm to MWTUI with `n`, switch text fields with `tab` and save the farm setting with `enter`.
- Select farm - After creating a new farm you'll be sent to the farm selection screen, otherwise press `a`. </br>
Navigate the menu with `↓ ↑`, select the desired farm with `enter` and wait for the miner list to refresh.
- Select miners - Use `enter` to multi select miners, press `tab` to see all selected miners. Miners can be deselected with `enter` as well.
- Send command - press the key associated with the desired command to send.
  - `w` - will wake machines.
  - `s` - will sleep machines.
  - `f` - will enable fastboot.
  - `s` - will disable fastboot.
  - `p` - will set pools - use `tab` to switch text fields and `enter` to send the command.
  - `l` - will set power limit - `enter` sends the new value.
- Deselect miners - At the moment there is no way to unselect all selected miners, to avoid issues deselect all selected miners.
- Refresh miner list - After sending a command and giving time for the command to process, press `m` to refresh the miner list and see the updated machine info.
- Commands can also be sent to a single machine by just having it highlighted and pressing a command key.

[comment]: <> (List of bugs and missing features)
### Known Bugs and Missing Features

- The view will break if your terminal window is too small.
- The view will break if resizing the terminal window.
- Setting pools has no validation check on both the pool URL and worker name.
- Setting farms has no IP validation.
- Setting power limit has no validation.
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
- When scanning for miners API errors can occur marking the machine as off and setting errors to `API_ERROR`.
  - The issue is related to the version of firmware on the whatsminer, currently unknown if its due to being too old or too new.

[comment]: <> (How to Install From Release)
### Installing from release

Download the appropriate release version from the release page and unzip it. </br>
Open a terminal and navigate to the newly created folder and run one of the following commands:</br>
> Windows: </br>
> `.\WMTUI.exe`

> Mac/Linux: </br>
> `./WMTUI`

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

