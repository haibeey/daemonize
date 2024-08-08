#  Daemonize
if you want to start a process in  background, Then you should checkout daemonize. Works on unix style systems.  

###  Usage  
daemonize --program python3 --args app.py --name pythonapp  
daemonize --program godoc --args -http:=6060  
daemonize --show
daemonize --name pythonapp  --kill.  
daemonize --kill --name test


###  Building  
go build in repo directory builds the binary

### example
```
+-------+------+-------+----------+---------+
|  PID  | NAME |  CPU  |  MEMORY  |  DISK   |
+-------+------+-------+----------+---------+
| 63085 | test | 12.9% | 154.4 MB | unknown |
+-------+------+-------+----------+---------+
```

## Download

```
wget -qO- https://raw.githubusercontent.com/haibeey/daemonize/master/install/install.sh | bash
```

#####  Contribution for creating services for os specific platform are welcome.
  
