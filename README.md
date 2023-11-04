# VFMP
My file management protocol (even though it is not a protocol and doesn't manage files)
# What can it do
## Count
Count all files in a directory (recursive)
## Index
Save paths to all files in a directory (also recursive) to a file on the drive
## Search
Searches for paths saved to the file.
### Regular search
Faster, useful when you know what you are looking for
### Fuzzy search
Slower, but useful when you can't remember the exact file name
# How to use it
## Starting the background process/daemon
### Linux
Run the `build/vfmpd` executable as root to start the daemon
### Other
Build the project in the `sservice/` directory and run with elevated permissions
## Interacting
### Using the command line
Build the project in the `cli/` directory, and run it.
You will get an interactive session connected to the VFMPd process.

Here is a list of commands and how to use them:
| Command | Arguments (positional) | What it does |
| --- | --- | --- |
| count | directory to count | Counts all the files in the provided directory and subdirectories |
| index | directory to index | Indexes the same files as count counts, but it saves them in a trie structure as a file on the drive |
| search | root directory (does nothing), search string (what to search), use fuzzy search (true/false) | Searches for a file in the trie. (If fuzzy search is used, returns wierd JSON object)
### Using the GUI
Build the project in the `gui/` directory, and run it.
You will get a desktop application with a basic UI to interact with
### Directly accessing the TCP server
The default port is 32768 on localhost, there is no documentation. If you want to use it, read the code, I tried to make it understandable