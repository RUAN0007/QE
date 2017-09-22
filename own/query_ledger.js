'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Hyperledger Fabric Sample Query Program
 */

var hfc = require('fabric-client');
var path = require('path');
const util = require('util')


var options = {
    wallet_path: path.join(__dirname, './creds'),
    user_id: 'PeerAdmin',
    channel_id: 'mychannel',
    chaincode_id: 'prov_cc',
    network_url: 'grpc://localhost:7051',
    txn_id: 'd54cc557cae8a9f8af0ec6ab7e80d9854456134d6354095c53fb57b294f59e31'
};

var channel = {};
var client = null;

Promise.resolve().then(() => {
    console.log("Create a client and set the wallet location");
    client = new hfc();
    return hfc.newDefaultKeyValueStore({ path: options.wallet_path });
}).then((wallet) => {
    console.log("Set wallet path, and associate user ", options.user_id, " with application");
    client.setStateStore(wallet);
    return client.getUserContext(options.user_id, true);
}).then((user) => {
    console.log("Check user is enrolled, and set a query URL in the network");
    if (user === undefined || user.isEnrolled() === false) {
        console.error("User not defined, or not enrolled - error");
    }
    channel = client.newChannel(options.channel_id);
    channel.addPeer(client.newPeer(options.network_url));
    return;
}).then(() => {
    console.log("Make query for a transaction");
    return channel.queryTransaction(options.txn_id)
}).then((processed_txn) => {
    console.log("returned from query");
    var validation_code = processed_txn.validationCode
    if(validation_code !== 0) {
      console.error("Invalid txn validation code", validation_code);
    }

    console.log(util.inspect(processed_txn, false, null))
}).catch((err) => {
    console.error("Caught Error", err);
});
