# gwire
An encrypted p2p chat system

## Features
- p2p chat, No 3rd party
- multi-user chat
- low on cpu and memory
- Encryption with AES
- Automatic encryption with RSA and AES

## Installation
Download the version for your operating system and run it. 
Optionally you can put it somewhere so you don't need to run the file locally from the directory it is in.
Other than that, there is no installation required

## Usage

### Flags
```
	-l, --listen				listen
	-p, --port 		[port]		use port [port]
	-h, --host 		[host]		Connect to [host]-(Ip)
	-v, --verbose		[level]		Show some more info levels: 1-3
	-u, --username 		[username]	set a username
	-t, --time			        enable timestamps
	-s, --slowmode		[seconds]	Enable slowmode
	-e, --encrypt		[password]	use AES instead of RSA
	-d, --no-encryption			Do not use encryption (Not recommended)
``` 
### Examples
`gwire -l -p 1234` listen on port 1234
<br>
`gwire -h 192.168.0.1 -p 1234` connect to 192.168.0.1 on port 1234 (encrypted)
<br><br>
`gwire -l -p 1234 -v -t -s 2 -u adzsx` listen on port 1234, use verbose mode, show timestamps, enable slowmode to 2 seconds and use adzsx as the username
<br><br>
`gwire -l -p 1234 -e Chk16QFV8xIj1tEyJhSjszkC5ERiAPwJ` Listen on port 1234, use manual encryption with AES and a specified password
<br><br>
`gwire -p 1234 -h 192.168.0.1 -d` Connect to 192.168.0.1 on port 1234, do not use encryption



## Contributing
Issues and Pull requests are welcome

# [License](https://choosealicense.com/licenses/gpl-3.0/)
