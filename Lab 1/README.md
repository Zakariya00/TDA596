**Setup**  
<1> Download files  
<2> Run console commands from file directory  
<3> If docker is not installed, install it.  
<4> Run Docker (should be running for next steps)  
  
**Build Commands (copy and paste)**  
<5> docker build --tag http_server .  
<6> docker build --tag proxy_server .  

**Run Commands (copy and paste)**  
<7> docker run -p 8080:8080 -t http_server  
    (can replace 8080:8080 with any ports as long as you change port args for server too)  
    
<8> Find url address for http_server in the format:  
    http://<ip-address:<portnumber>    
<9> Change proxy_server Dockerfile args (arg 2) to the url address you found  
  
<10> docker run -p 1234:1234 -t proxy_server  
    (can replace 1234:1234 with any ports as long as you change port args for server too)  
  
<11> Servers should be running now  

Command-line-arguments are supplied in the last line of the Dockerfile  
ENTRYPOINT ["/http_server", "arg1", "arg2"]  

Can change port for server and/also mainserver for proxy_server  
On code change rerun commands, 5 to 10.  

If you change the ports in Dockerfile then change them in docker run commands also:  
docker run -p serverport:serverport -t http_server  
