  **Cloud Acess**  
  Server(s) available on AWS Cloud  (current addresses):  
  http_server URL  : http://ec2-18-209-56-42.compute-1.amazonaws.com:8080/  
  proxy_server URL : http://ec2-18-209-56-42.compute-1.amazonaws.com:1234/  
      
  
  GET Requests: Possible both through the browser and command promt  
  POST Requests: Works with supplied html file named index (change sending address on need).  
  Has been uploaded to cloud server    
  Link: http://ec2-18-209-56-42.compute-1.amazonaws.com:8080/uploadMenu.html 
    
  Docker Repositories (Public):  
  http_server   : https://hub.docker.com/repository/docker/zakariya00/http_server  
  proxy_server  : https://hub.docker.com/repository/docker/zakariya00/proxy_server 

 
    

------------------------------------------------------------------------
***Setup***  
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
  
------------------------------------------------------------------------
  **Solution Explaination**
    
    
  http_server & proxy_server listen on the port given as command-line-arguments.  
  http_server can run wwithout an arg, by using default port ":8080".  
  proxy_server can not run without args, needs port and forwarding address.  
    
    
  http_server employs limits the number of clients being served by using semaphores.  
  Semaphores are hard-coded but easily configured to suppourt serving more client.  
  Semaphore blocks and waits when 0. Avoiding busy wait that drains resources.  
    
    
  http.serve runs in a loop and automatically starts a new go process when a http request comes in.  
  Therefore can not be blocked with a semaphore before spawning a new go process. So requires to block and wait in the spawned handler. 
  As such we still suffer the extra resource costs of spawning a new go process but they remain blocked until they acquire a semaphore.  
  http_server suppourts methods GET and POST, sending error message on the rest.  
    
  proxy_server handles only GET requests. Making a new GET request to the main server, waiting for the response then copying it in  
  the original client request response to the proxy_server. Handles any other method requests appropriately.
