# Estuary Barge
Barge is a CLI tool that directly uses estuary api to stream/upload files to the filecoin network.

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

### Configuration file
```
{
  "estuary": {
    "host": "http://localhost:3004",
    "primaryshuttle": "http://localhost:3005",
    "token": "<local API token>"
  }
}
```