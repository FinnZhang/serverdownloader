# serverdownloader

With the help of [echo](https://github.com/labstack/echo) Go web framework. I can make this tool easily.

## Usage
* Download code  
``` go get github.com/FinnZhang/serverdownloader ```    
* Build  
``` go build ./serverdownloader.go ```  
* Run  
index.html must be in the same directoty with the serverdownloader binary. <br/>(If you use `go install` or defalut `go get` command,Ensure that two files in the same directory.)Run the binary and open a web browser,the location is your vps ip.
## Tips   
1. Install `tmux` by `sudo apt install tmux`,Then run `tmux` command.
2. Run the `serverdownloader` bianry file in tmux.
3.  `Ctrl+b`  then  `d` leave it running it in the background.

## Other
This tool is listen to 80 port.Maybe used already or need to be root to listen to 80.you can change the port and then rebuild it.


