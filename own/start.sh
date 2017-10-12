#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

if [ ! -d ~/.hfc-key-store/ ]; then
	mkdir ~/.hfc-key-store/
fi
cp $PWD/creds/* ~/.hfc-key-store/
# launch network; create channel and join peer to channel
cd ../basic-network
./start.sh

function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

removeUnwantedImages 
# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
docker-compose -f ./docker-compose.yml up -d cli

# Install example chain code query
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n supplychain -v 1.0 -p github.com/supplychain

# Init the material
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n supplychain -v 1.0 -c '{"Args":["init","1","1","1", "1", "2", "1","1","1","DBS", "1000"]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
sleep 05

echo "=========================Making Camera=========================="

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["MakeCamera","FrontCam0","BackCam0","Camera0"]}'
sleep 05

# Make CPU
echo "=========================Making CPU=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["MakeCPU","ALU0","ControlUnit0","Register0", "Register1", "CPU0"]}'
sleep 05

# Make Mainboard
echo "=========================Making Mainboard=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["MakeMainboard","CPU0","Memory0","SSD0", "Mainboard0"]}'
sleep 05

# Assemble Iphone
echo "=========================Assemble IPhone=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Assemble","Camera0","Battery0","Mainboard0", "IPhone0", "Manufacturer0"]}'
sleep 05

# Procure Iphone to Retailer from  manufacuturer
echo "=========================Procure IPhone=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Procure","IPhone0","Manufacturer0","Retailer0"]}'
sleep 05

# Purchase Iphone to Retailer from retailer
echo "=========================Purchase IPhone=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Purchase","IPhone0","Customer0","DBS", "Retailer0", "100"]}'
sleep 05

# Resell Iphone to Retailer from retailer
echo "=========================Resell IPhone=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Resell","IPhone0","Customer0","DBS", "Customer1", "50"]}'


# Get the bank account
# echo "=========================Query Bank Account=========================="
# docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Query","DBS"]}'
# sleep 05

# Get the IPhone
echo "=========================Query IPhone=========================="
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n supplychain -c '{"Args":["Query","IPhone0"]}'
# sleep 05

# printf "\nTotal execution time : $(($(date +%s) - starttime)) secs ...\n\n"
# STATEDB_SIZE="$(docker exec peer0.org1.example.com  du -s --block-size=K  /var/hyperledger/production/ledgersData/stateLeveldb/ | sed 's/[^0-9]//g')"
# echo "State_DB SIZE: $STATEDB_SIZE"
# HISTORYDB_SIZE="$(docker exec peer0.org1.example.com  du -s --block-size=K  /var/hyperledger/production/ledgersData/historyLeveldb/ | sed 's/[^0-9]//g')"
# echo "History_DB SIZE: $HISTORYDB_SIZE"

