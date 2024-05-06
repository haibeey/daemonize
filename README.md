#  Daemonize
if you want to start a process in  background, Then you should checkout daemonize  

###  Usage  
daemonize -program python3 -args app.py -name pythonapp  
daemonize -program godoc -args -http:=6060  
daemonize -show 1 
daemonize -name pythonapp  -kill 1 


###  Building  
go build in repo directory builds the binary

#####  Contribution for creating services for os specific platform are welcome.
