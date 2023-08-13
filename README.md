<p align="center">
  <img src="https://github.com/adzsx/gwire/blob/main/img/gwire-banner.png">
</p>

## Features
- p2p chat system
- multi-user chat
- low on cpu and memory
- Encryption with AES
- Automatic encryption with RSA and AES

## Installation
Download the version for your operating system and run it. 
Optionally you can put it somewhere so you don't need to run the file locally from the directory it is in.
Other than that, there is no installation required

## Usage

### Flags2
```
    -l, --listen			        listen
    -p, --port 	    [port]	        use port [port]
    -h, --host 	    [host]	        Connect to [host]-(Ip)
    -v, --verbose		            Show some more info
    -u, --username 	[username]	    is didsplayed for other users
    -t, --time			            enable timestamps
    -s, --slowmode	[seconds]	    Enable slowmode, minimum is 0.1s
    -e, --encrypt	([password])    If [password] is given, use AES, if not, encrypt automatic with RSA
```
### Examples
`gwire -l -p 1234` listen on port 1234
<br>
`gwire -h 192.168.0.1 -p 1234` connect to 192.168.0.1 on port 1234
<br><br>
`gwire -l -p 1234 -v -t -s 2 -u adzsx` listen on port 1234, use verbose mode, show timestamps, enable slowmode to 2 seconds and use adzsx as the username
<br><br>
`gwire -l -p 1234 -e Chk16QFV8xIj1tEyJhSjszkC5ERiAPwJ` Listen on port 1234, use manual encryption with AES and a specified password
<br><br>
`gwire -p 1234 -h 192.168.0.1 -e` Connect to 192.168.0.1 on port 1234, use automatic encryption


## Contributing
Issues and Pull requests are welcome

# [License](https://choosealicense.com/licenses/gpl-3.0/)
