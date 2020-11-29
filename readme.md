Mutagen Monitor for Mac
=======================

Description
-----------
[Mutagen Monitor](https://github.com/andrewmed/mutagenmon) is a system bar only Mac OSX application, made for monitoring active [Mutagen](https://mutagen.io) sessions

Ratio
-----
I have been using Mutagen for quite a lot for developing. For this purpose I use it only in 
[one-way-safe synchronization mode](https://mutagen.io/documentation/synchronization). I found Mutagen to be very convenient for syncing large monorepos because my main production environment is Linux and for editing I use Mac. In backend we have clear separation between binary files that are symlinks mostly and the code. Mutagen fit in nicely with a one-way-safe mode and ignore symlinks flag.

The only thing I was missing was a monitoring agent for Mac. That is why I made this

How to use
----------
![Image](demo.png)

In the bar there is a number of sessions in these states:
* healthy (connected and with no file conflict)
* connected (regardless of file conflicts)
* registered (regardless of file conflicts of connection state)
 
So on the picture we have one session, connected, but it has file conflicts. On mouse over you can see details for the session

To synchronize beta (remote) with alpha (local) just click on menu item with corresponding session. Conflicting files/directories on beta will be overwritten.

How to build
------------
```
./build.sh
```

Alternatives
----------
[Mutagen Monitor](https://github.com/andrewmed/mutagenmon) was inspired by [MutagenMon](https://github.com/rualark/MutagenMon) thanks to @rualark! But that one seems to be aiming Windows, and there is also "python vs go" difference

This monitor is written in Go and uses native api (Mutagen itself is written in Go), so no new extra process is created while monitor is watching. This makes it possible to have less than 0.0% CPU usage when monitor is in background

FAQ
----
Q: I run Mutagen Monitor, but nothing happens, there is no menu in system bar

A: Check if Mutagen [daemon](https://mutagen.io/documentation/introduction/daemon) is running. To start daemon just do: `mutagen list`

Licence
-------
GPLv3