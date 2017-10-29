# GetToDahChoppa
An executable to generate multiple sub-jobs from indicated GCODE file, suitable for printing color-by-layer effects.

## Overview
This initial version works with the Cura slicer software. In subsequent versions, I may add support for other softwares but I would need access to GCODE examples first.

There are other methods of printing 3D parts with different colors in vertical patterns. I have found these to often be problematic, leading to a hot extruder crashing into the part, for example. In another case, running the filament change wizard appeared *not* to reset the extruder positioning mode back to what it was before, resulting in the hot extruder drilling into my print bed.

So I wanted to break up the original GCODE file into multiple ones. In theory, this should allow us to then send each job to the printer in order and to prep the printer in each case by running the appropriate filament change activity.

If you're good with GCODE, you can then further edit each sub-job file to make sure that it's not doing some odd extruder assembly move which would take it through the path of the partially-printed part sitting on the print bed.

Although untested at the time of this writing, you should be able to print different filament types per sub-job given that you can now edit each file for the temperature settings at the top.

## Syntax
This assumes that you've installed it per the instructions below so that it will be in your path.

```
# Example of the first file of two
$ GetToDahChoppah -beginning -to 12 -ordinal 1 -count 2 -msg="First File" PathToGcodeFile
Original:   /Users/user/Desktop/filename.gcode
Slicer:     Cura_SteamEngine 2.3.1
Layers:     24
Slicing:
  Output filename:  /Users/user/Desktop/filename_1of2.gcode
  From:             0
  To:               12
  Msg:              "First File"
  
# Example of the second file of two
$ GetToDahChoppah -from 13 -end -ordinal 2 -count 2 -msg="Second File" PathToGcodeFile
Original:   /Users/user/Desktop/filename.gcode
Slicer:     Cura_SteamEngine 2.3.1
Layers:     24
Slicing:
  Output filename:  /Users/user/Desktop/filename_2of2.gcode
  From:             13
  To:               24
  Msg:              "Second File"
```

## Installation
The installation of this depends upon whether or not you also have a computer based upon the OS X operating system and further, whether or not you already have the Go compiler itself.

### Mac
Here are the instructions for installing this executable if you are on an Apple-based computer.

#### If you have Go installed:

```
$ cd /usr/local/go/bin
$ sudo curl https://github.com/OutsourcedGuru/GetToDahChoppa/raw/master/bin/GetToDahChoppa GetToDahChoppa
$ cd ~/Desktop
$ which GetToDahChoppa
/usr/local/go/bin/GetToDahChoppa
$ GetToDahChoppah -info PathToGcodeFile
Original:   /Users/user/Desktop/filename.gcode
Slicer:     Cura_SteamEngine 2.3.1
Layers:     24

Finished.
```

#### If you have don't have Go installed:

```
$ cd /usr/local/bin
$ sudo curl https://github.com/OutsourcedGuru/GetToDahChoppa/raw/master/bin/GetToDahChoppa GetToDahChoppa
$ cd ~/Desktop
$ which GetToDahChoppa
/usr/local/GetToDahChoppa
$ GetToDahChoppah -help
Usage of GetToDahChoppa:
  -beginning
    	start from first layer seen in file
  -count int
    	job count [-count=3] with "_1of3" appended to filename as output (default -1)
  -end
    	end at last layer seen in file
  -from int
    	start from indicated layer [-from=4] (default -1)
  -info
    	just read layer data from file
  -msg string
    	M117 message info to add to GCODE [-msg="Blue PLA"] (default "N/A")
  -ordinal int
    	job ordinal [-ordinal=1] (default -1)
  -to int
    	end at indicated layer [-to=24] (default -1)
```

### Windows or UNIX
Here are the instructions for building this executable if you are on a different operating system.

#### Install Go
The first step is to [install the Go language compiler](https://golang.org).

It's then usual to create a Go working folder under your user's profile.

```
# These two are optional, depending upon whether
# or not you did this during the Go installation
$ mkdir -p ~/go/src
$ cd ~/go/src
$ go get github.com/OutsourcedGuru/GetToDahChoppa/
```

This should download everything required and build it for you. Assuming that you installed Go correctly earlier and it's in your path, you should then be able to run it as in the instructions above.

## Important Notes
* The current version expects Cura as the slicer and has been tested with v2.3.1 of same. It has been tested with the Robo C2 printer with a single extruder and no heated bed.
* Given the behavior of G29 autoleveling routines it is unwise to allow autoleveling on any of the sub-jobs because the nine-point testing should interfere with already-printed parts on the bed itself. For this reason I would strongly suggest toggling off autoleveling, per [OctoPrint-plugin-toggle-autolevel](https://github.com/OutsourcedGuru/OctoPrint-plugin-toggle-autolevel).
* Additionally, it may be worth noting the current path(s) used during the printing of a priming line in the front of the printer so that the extruder assembly doesn't crash into your part. Proper placement of your part during slicing may be the best course of action here. Lastly, I would suggest that you remove each priming line as printed since each sub-job will attempt to recreate it.
* You might consider installing [M117Popup](http://plugins.octoprint.org/plugins/M117PopUp/), [M117NavBar](http://plugins.octoprint.org/plugins/M117NavBar/) or [StatusLine](http://plugins.octoprint.org/plugins/status_line/) as well as [DisplayZ](http://plugins.octoprint.org/plugins/displayz/) while you're at it to take advantage of the `-msg` flag option of GetToDahChoppa.