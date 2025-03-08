None of the tools do what I want, nor can I actualy get python working on my computer
Python sucks, but C & C++ suck more

`/pyboard` is a go implementation of the pyboard.py class

This tool will have two main functions
dump & sync

dump: make a copy of the entire root dir of a micropython device to a local folder of your choosing
sync: keep the microcontroller in sync with all file changes made on the local device and restart the program when changes occor

I may also add simple fs commands and or a shitty bash equivalent 

This tool is not supposed to be a replacement for current tools, I am writing this purely for myself
You are not supposed to use this, I will not provide support for this, but if you with to make a PR I will review it and merge it given it works