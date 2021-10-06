# Car Audio Database

This project endeavors to become a central repository for storing frequency
response characteristics of speakers when paired with specific vehicles in 
specific positions.

# Development

## Angular

From the `car-audio-database` folder, run:
```shell
npm install
ng serve
```

Point your browser to http://localhost:4200

## GoLang
From the `server` folder, run:
```shell
go run *.go
```

### Testing server via the Angular application
With the server running, run this command from the `car-audio-database` folder:
```shell
ng build
```

This will re-generate the files under the `dist` directory, from which the Go 
server's Angular files are served.

Point your browser to http://localhost:8080

#### Why not `ng serve`?
Modern browsers will block CORS requests between different ports on localhost

But, if the Angular app is running on the same port as the server, it will 
allow sensitive requests (like file uploads).

## Docker 

### Build
To build a new image, run from the project root:
```shell
docker build . --tag bradj.ca/car-audio-db
```

### Run
To run the server via docker, run from the project root
```shell
docker run -p 8080:8080 bradj.ca/car-audio-db
```