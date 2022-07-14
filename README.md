# Estuary Barge
Barge is a CLI tool that directly uses estuary api to stream/upload files to the filecoin network.

![image](https://user-images.githubusercontent.com/4479171/178054186-70c482f9-679d-4ab0-9f3d-e8ffa6ce49a7.png)


*This project is based on the a sub module on Estuary under /cmd/barge.* 

## Installation
### Pre-requisites
- Go
- make
- filecoin-ffi (this is a submodule)

Clone this repo and run the following
```
make all
```
This will generate a `barge` binary on the root folder that you can test.

## Usage

### Grab your API key from Estuary and run the following:
```
./barge login <API KEY>
```

### Initialize
Initialize barge with the following command. This will create a configuration file which
holds the estuary connection information.
```
./barge init 
```

### Local configuration file
```
{
  "estuary": {
    "host": "http://localhost:3004",
    "primaryshuttle": "http://localhost:3005",
    "token": "<local API token>"
  }
}
```

### Remote Estuary configuration
```
{
  "estuary": {
    "host": "https://api.estuary.tech",
    "primaryshuttle": "https://shuttle-4.estuary.tech",
    "token": "<Estuary API token>"
  }
}
```

## Usage
### Upload a file
```
./barge plumb put-file <file path>
```

### Upload a CAR file
```
./barge plumb put-car <CAR file path>
```

## WebUI
Run Web
```
./barge web
```

[http://localhost:3000](http://localhost:3000)

## REST API Endpoints
Run Web and use the following endpoints to interact with the barge.
```
./barge web
```

### Upload file
```
curl --location --request POST 'http://localhost:3000/api/v0/plumb/file' \
--form 'file=@"website.png"'
```

### Upload CAR
```
curl --location --request POST 'http://localhost:3000/api/v0/plumb/car' \
--form 'file=@"file.car"'
```