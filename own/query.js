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
    chaincode_id: 'supplychain',
    network_url: 'grpc://localhost:7051',
};

var channel = {};
var client = null;

// Return the promise of provenance records and 
//   dependent read version of the asset of a specific version

// a list of two objects
//   the first is the provenance records
//   the second is the list of dependent reads and its version
function GetDependency(asset, block_num, tx_num) {
    return Promise.resolve().then(() => {
        // console.log("Create a client and set the wallet location");
        client = new hfc();
        return hfc.newDefaultKeyValueStore({ path: options.wallet_path });
    }).then((wallet) => {
        // console.log("Set wallet path, and associate user ", options.user_id, " with application");
        client.setStateStore(wallet);
        return client.getUserContext(options.user_id, true);
    }).then((user) => {
        // console.log("Check user is enrolled, and set a query URL in the network");
        if (user === undefined || user.isEnrolled() === false) {
            throw new Error("User not defined, or not enrolled - error");
        }
        channel = client.newChannel(options.channel_id);
        channel.addPeer(client.newPeer(options.network_url));
        return;
    }).then(() => {
        return channel.queryBlock(block_num);
    }).then((block) => {
        // console.log("Retrieved txn: ");
        // console.log(util.inspect(block.data.data[0], false, null));

        var txn_num = block.data.data.length;
        // console.log("Getting Block: ", block, " with ", txn_num, " txns");
        if (txn_num <= tx_num) {
            throw new Error("Total number of transactions is less than the transaction idx");
        }

        // console.log("Retrieved txn: ");
        // console.log(util.inspect((block, false, null)));

        var actions = block.data.data[tx_num].payload.data.actions;
        if(actions.length !== 1) {
          throw new Error("Multiple action detected: ", actions.length);
        }

        // console.log("Retrieved action: ");
        // console.log(util.inspect(actions[0], false, null));

        var ns_rwsets = actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset

        var cc_rwset; // Read-write set for this chain code

        for (var rwset_id in ns_rwsets) {
            if (ns_rwsets[rwset_id].namespace === options.chaincode_id) {
               cc_rwset = ns_rwsets[rwset_id].rwset;
               break; 
            }
        }

        if (cc_rwset === undefined) {
            throw new Error("Chaincode " + options.chaincode_id + " does not enable provenance tracking. "); 
        } 

        var cc_readset = cc_rwset.reads;
        var cc_wrtset = cc_rwset.writes;

        // console.log(cc_wrtset);

        var provenance; // provenance records for the asset
        for (var wrt_id in cc_wrtset) {
            if (cc_wrtset[wrt_id].key == asset + "_prov") {
              provenance = JSON.parse(cc_wrtset[wrt_id].value);
              break; 
            }
        }

        if (provenance === undefined) {
          throw new Error("The provenance for asset " + asset + " is not found"); 
        } 


        // Start to get the version of the dependency reads
        var result = []
        var dependency_read_versions = []
        for (var dep_read_idx in provenance.DepReads) {
            var dep_read_asset = provenance.DepReads[dep_read_idx];
            for (var read_idx in cc_readset) {
                if (cc_readset[read_idx].key == dep_read_asset) {
                    dependency_read_versions.push(cc_readset[read_idx]);
                }
            }  // end for
        }  // end for

        return [provenance, dependency_read_versions];
    });
}


// Return the promise of provenance records and dependent read version of the latest asset
// a list of two objects
//   the first is the provenance records
//   the second is the list of dependent reads and its version
function GetLastestDependency(asset) {
    return Promise.resolve().then(() => {
        // console.log("Create a client and set the wallet location");
        client = new hfc();
        return hfc.newDefaultKeyValueStore({ path: options.wallet_path });
    }).then((wallet) => {
        // console.log("Set wallet path, and associate user ", options.user_id, " with application");
        client.setStateStore(wallet);
        return client.getUserContext(options.user_id, true);
    }).then((user) => {
        // console.log("Check user is enrolled, and set a query URL in the network");
        if (user === undefined || user.isEnrolled() === false) {
            throw new Error("User not defined, or not enrolled - error");
        }
        channel = client.newChannel(options.channel_id);
        channel.addPeer(client.newPeer(options.network_url));
        return;
    }).then(() => {
        // console.log("Make query");
        var transaction_id = client.newTransactionID();
        // console.log("Assigning transaction_id: ", transaction_id._transaction_id);

        const request = {
            chaincodeId: options.chaincode_id,
            txId: transaction_id,
            fcn: 'latest_txn',
            args: [asset]
        };
        return channel.queryByChaincode(request);
    }).then((query_responses) => {
        if (query_responses[0] instanceof Error) {
            throw new Error("error from query = ", query_responses[0]);
        }
        console.log("returned from query");
        if (!query_responses.length) {
            console.error("No payloads were returned from query");
        } 

        if (query_responses.length > 1) {
            console.error("Only single payload is required from the query.")
        } 

        // console.log("Latest Written tnx ID for ", asset, 
        //             " is ", query_responses[0].toString());
        return query_responses[0].toString()
    }).then((last_wrt_txn_id) => {
        // console.log("Make query for a transaction ", last_wrt_txn_id);
        return channel.queryTransaction(last_wrt_txn_id);
    }).then((processed_txn) => {
        console.log("returned from query");
        var validation_code = processed_txn.validationCode
        if(validation_code !== 0) {
          throw new Error("Invalid txn validation code", validation_code);
        }

        var actions = processed_txn.transactionEnvelope.payload.data.actions
        if(actions.length !== 1) {
          throw new Error("Multiple action detected: ", actions.length);
        }


        var ns_rwsets = actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset
        console.log("# of ns_rwsets: ", ns_rwsets.length)

        var cc_rwset; // Read-write set for this chain code

        for (var rwset_id in ns_rwsets) {
            if (ns_rwsets[rwset_id].namespace === options.chaincode_id) {
               cc_rwset = ns_rwsets[rwset_id].rwset;
               break; 
            }
        }

        if (cc_rwset === undefined) {
            throw new Error("Chaincode " + options.chaincode_id + " does not enable provenance tracking. "); 
        } 


        var cc_readset = cc_rwset.reads;
        var cc_wrtset = cc_rwset.writes;

        var provenance; // provenance records for the asset
        for (var wrt_id in cc_wrtset) {
            if (cc_wrtset[wrt_id].key == asset + "_prov") {
              provenance = JSON.parse(cc_wrtset[wrt_id].value);
              break; 
            }
        }

        if (provenance === undefined) {
          throw new Error("The provenance for asset " + asset + " is not found"); 
        } 

        // Start to get the version of the dependency reads
        var dependency_read_versions = []
        for (var dep_read_idx in provenance.DepReads) {
            var dep_read_asset = provenance.DepReads[dep_read_idx];
            for (var read_idx in cc_readset) {
                if (cc_readset[read_idx].key == dep_read_asset) {
                    dependency_read_versions.push(cc_readset[read_idx]);
                }
            }  // end for
        }  // end for

        return [provenance, dependency_read_versions];
    });
};

// GetLastestDependency('IPhone0').then((provenance) => {
//     console.log(util.inspect(provenance, false, null));
//     // console.log("Txn_NUmber for IPhone0: ", provenance[1][0]["version"]["tx_num"].toInt());
//     // console.log("Block Number for IPhone0: ", provenance[1][0]["version"]["block_num"].toInt());
// }).catch((err) => {
//     console.error("Caught Error", err);
// });


// GetDependency('IPhone0', 6, 0).then((provenance) => {
//      console.log("Dependency for IPhone0: ", provenance);
//  }).catch((err) => {
//     console.error("Caught Error", err);
// });


// Trace Line IPhone -> IPhone -> IPhone -> Mainboard -> CPU -> ALU

console.time('level0');
console.time('level1');
console.time('level2');
console.time('level3');
console.time('level4');
console.time('level5');
console.time('level6');

GetLastestDependency('IPhone0').then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");
    var dep_reads = result[1];
    var i, pre_blk_num, pre_txn_num;
    for (i = 0; i < dep_reads.length; ++i) {
      if(dep_reads[i]["key"] !== "IPhone0") {
        console.log("Add to Account: ", dep_reads[i]["key"]);
      } else {
        pre_blk_num = dep_reads[i]["version"]["block_num"].toInt();
        pre_txn_num = dep_reads[i]["version"]["tx_num"].toInt();
      }
    }
    console.log("=======================================================");
    console.timeEnd('level0');
    return GetDependency('IPhone0', pre_blk_num, pre_txn_num); 

}).then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");
    var dep_reads = result[1];
    var i, pre_blk_num, pre_txn_num;
    for (i = 0; i < dep_reads.length; ++i) {
      if(dep_reads[i]["key"] !== "IPhone0") {
        console.log("Deduct From Account: ", dep_reads[i]["key"]);
      } else {
        pre_blk_num = dep_reads[i]["version"]["block_num"].toInt();
        pre_txn_num = dep_reads[i]["version"]["tx_num"].toInt();
      }
    }
    console.log("=======================================================");
    console.timeEnd('level1');
    return GetDependency('IPhone0', pre_blk_num, pre_txn_num); 

}).then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");


    var dep_reads = result[1];
    var pre_blk_num = dep_reads[0]["version"]["block_num"].toInt();
    var pre_txn_num = dep_reads[0]["version"]["tx_num"].toInt();
    console.log("=======================================================");

    console.timeEnd('level2');
    return GetDependency('IPhone0', pre_blk_num, pre_txn_num); 
}).then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");

    var dep_reads = result[1];
    var i, pre_blk_num, pre_txn_num;
    console.log("Dependent Components: ");
    for (i = 0; i < dep_reads.length; ++i) {
      console.log("  ", dep_reads[i]["key"]);
      if(dep_reads[i]["key"] === "Mainboard0") {
        pre_blk_num = dep_reads[i]["version"]["block_num"].toInt();
        pre_txn_num = dep_reads[i]["version"]["tx_num"].toInt();
      }    
    }
    console.log("=======================================================");
    console.timeEnd('level3');
    return GetDependency('Mainboard0', pre_blk_num, pre_txn_num); 
}).then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");

    var dep_reads = result[1];
    var i, pre_blk_num, pre_txn_num;

    console.log("Dependent Components: ");
    for (i = 0; i < dep_reads.length; ++i) {
      console.log("  ", dep_reads[i]["key"]);
      if(dep_reads[i]["key"] === "CPU0") {
        pre_blk_num = dep_reads[i]["version"]["block_num"].toInt();
        pre_txn_num = dep_reads[i]["version"]["tx_num"].toInt();
      }    
    }
    console.log("=======================================================");
    console.timeEnd('level4');
    return GetDependency('CPU0', pre_blk_num, pre_txn_num); 
}).then((result) => {
    var prov = result[0];
    var func_name = prov["FuncName"];
    console.log("===================", func_name, "=====================");

    var dep_reads = result[1];
    var i;

    console.log("Dependent Components: ");
    for (i = 0; i < dep_reads.length; ++i) {
      console.log("  ", dep_reads[i]["key"]);
    }
    console.log("=======================================================");
    console.timeEnd('level5'); 

}).catch((err) => {
    console.error("Caught Error", err);
});
