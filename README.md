# Prerequiste
  * Install relevant [Hyperledger Fabric 1.0 Prerequisite](http://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)
  * Download the forked Hyperledger Fabric 1.0 [fork](https://github.com/RUAN0007/fabric) to build the relevant docker images. 

## Query Latency
First Setup the network via the command
```
cd own/; ./start.sh
```

Run the NodeJS SDK query to get latency for provenance level. 
```
node query.js
```

## Query Storage Size by varying data size
```
NUM_IPHONE=50; ./workload.sh
```

## Graph Plotting
* Install [pyplot](https://matplotlib.org/api/pyplot_api.html) 
* Refer to python scripts in own/plot
