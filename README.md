
## About

Here we Go again 🔥. 

Once upon a time, after work, I dived deeper into how the Go runtime works. 

I was reading about Go runtime package and decided to make a lib for monitoring a system and send different metrics to a server. After creating an agent to collect metrics, I decided to take it a step further and build a server. I also wanted to explore the Go standard library and improve my skills with it. That's how this repository was created.  


## How to launch 

### Build step:
- As usual clone the repository
- From the root repository dir run
```
 - make build
 ```
### Launch step
#### 1. Launch using RAM storage
To run the system, simply start the server and agent in separate terminals without any flags or environment variables.
```
 - cd build
 - ./server (optionally --help) or ./agent (optionally --help)
 ```

#### 2. Launch using Docker
To run the system using docker compose, create .env file before. Run from the root dir:
```
 - cp .env.example .env
 - sudo docker compose build
 - sudo docker compose up
 ```
All environment variables are properly set for use with Docker Compose.

### Explore step

To make sure the system works correctly check logs or ping endpoint:
```
 http://$SERVER_ADDR:$SERVER_PORT/ping or localhost:8080/ping
 ```

To gather metrics, you can use various methods. Begin by using curl or a browser to access the root router:
```
 localhost:8080/
 ```

### Used libraries
 - github.com/gin-gonic/gin
 - github.com/go-resty/resty
 - github.com/jackc/pgx
 - github.com/stretchr/testify
 - github.com/google/uuid
