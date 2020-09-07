Mutagen Monitor for Mac
=======================

Description
-----------
[Mutagen Monitor](https://github.com/andrewmed/mutagenmon) is a system-bar-only Mac OSX application, made for monitoring active [Mutagen](https://mutagen.io) sessions

I have been using Mutagen actively for quite a lot for developing. For this purpouse I use it only in 
[one-way-safe synchronization mode](https://mutagen.io/documentation/synchronization).
I found Mutagen to be very convenient for syncing large monorepos because my main production environment is Linux and for editing I use Mac

The only thing I was missing was a monitoring agent for Mac. So I made this one

How to use
----------
![Image](demo.png)

In the bar there is a number of healthy vs total sessions (1/2)

Healthy are "staging", "watching" and "waiting" states provided that the session has no file conflicts. Any other state is unhealthy. So on the picture we have two sessions and only one of them is healthy. On mouse over you can see details for a session

There is no interaction yet through the menu, except for Quit (see **todo** below for more)

How to build
------------
```
git clone ...
cd mutagenmon/mutagenmon
./build.sh
```

There is also a prebuilt Mac OSX [application package](https://github.com/andrewmed/mutagenmon/releases), see [Releases](https://github.com/andrewmed/mutagenmon/releases) for it.


Alternatives
----------
This piece of software was inspired by [MutagenMon](https://github.com/rualark/MutagenMon) thanks to @rualark!
But that one seemed to be focused on Windows, and there is also "python vs go" difference.

This monitor is written in Go and uses native Mutagen api (it is written in Go also), so no process calls are made by the monitor.
This makes it possible to have less than 0.0% CPU usage when monitor is in background (on my machine)

MutagenMon by @rulark has more functionality however (conflict resolution feature is very nice).

Todo
----
* add conflict resolution (remove conflicting files on beta) on menu click (will work only for [one-way-safe synchronization mode](https://mutagen.io/documentation/synchronization))

Licence
-------
MIT