#  Daemonize
if you want to start a process in  background, Then you should checkout daemonize  

###  Usage  
daemonize -f python3 app.py  
daemonize python3 app.py  
deamonize -s -f python3 app.py  
deamonize -s python3 app.py  
 ./daemonize -f /Users/haibeey -b python3 -a flasktest.py  
 daemonize -b godoc -a -http:=6060  

The -s parameter simply tells daemonize to create a system service for you platform and start you applications  
The -f parameter exist to tell daemonize which folder should be consider as working directory  
daemonize -b godoc -o ch.out -a -http=localhost:6060  
#### note: all args passed to binary should be after any params  

###  Building  
go build in repo directory builds the binary

#####  Contribution for creating services for os specific platform are welcome.
