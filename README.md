
## About

Here we Go again 🔥. 

Once upon a time, after work, I dived deeper into how the Go runtime works. 

I was reading about Go runtime package and decided to make a lib for monitoring a system and send different metrics to a server. After creating an agent to collect metrics, I decided to take it a step further and build a server. I also wanted to explore the Go standard library and improve my skills with it. That's how this repository was created.  


## How to launch 

### Build step:
- As usual clone the repository
- Change the workdir to the repository dir and run
```
 - make build
 ```
### Launch step
To run the system, simply start the server and agent in separate terminals without any flags or environment variables.
```
 - cd build
 - ./server (optionally --help) or ./agent (optionally --help)
 ```

To make sure the system works correctly check logs.

To gather metrics, you can use various methods. Begin by using curl or a browser to access the root router (http://$SERVER_ADDR:$SERVER_PORT/)