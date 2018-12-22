# serverdownloader

1. Download file from source to server then to client.
2. Upload file from local to server.

## Usage
* Download code  
``` go get github.com/FinnZhang/serverdownloader ```    
* Build  
``` go build ./serverdownloader.go ```  
* Run  
```./serverdownloader```   
If you use `go install` or defalut `go get` command,Ensure that html files and binary are in the same directory.Run the binary and open a web browser,the address is ```http://{your_vps_ip}:{port}```.
## Tips   
1. Install `tmux` by `sudo apt install tmux`,Then run `tmux` command.
2. Run `./serverdownloader` in tmux.
3.  `Ctrl+b`  then  `d` leave it running in the background.

## Other
This tool is listen to 80 port.Need to be root to run it.You can change the listening port and then rebuild it.


