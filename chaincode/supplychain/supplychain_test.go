/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkIPhoneOwner(t *testing.T, stub *shim.MockStub, serial string, expected_owner string) {
	iphone_bytes := stub.State[serial]
	if iphone_bytes == nil {
		fmt.Println("Entity ", serial, "doesn't exist. ")
		t.FailNow()
	}

	var iphone Iphone

	err := json.Unmarshal(iphone_bytes, &iphone)
	if err != nil {
		fmt.Println("Fail to unmarshal entity with ID ", serial)
		t.FailNow()
	}

	if iphone.Owner != expected_owner {
		fmt.Println("Iphone ", serial, " is owned by ", iphone.Owner, " NOT ", expected_owner)
		t.FailNow()
	}
}

func checkEntityUsage(t *testing.T, stub *shim.MockStub, serial string, expected_used bool) {
	entity_bytes := stub.State[serial]
	if entity_bytes == nil {
		fmt.Println("Entity ", serial, "doesn't exist. ")
		t.FailNow()
	}

	var actual_entity Entity

	err := json.Unmarshal(entity_bytes, &actual_entity)
	if err != nil {
		fmt.Println("Fail to unmarshal entity with ID ", serial)
		t.FailNow()
	}

	if actual_entity.Used != expected_used {
		if actual_entity.Used {
			fmt.Println("Entity ", serial, " is already used. ")
		} else {
			fmt.Println("Entity ", serial, " is not used. ")
		}
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func Test(t *testing.T) {
	scc := new(SupplyChaincode)
	stub := shim.NewMockStub("ex02", scc)

	// Init A=123 B=234
	checkInit(t, stub, [][]byte{[]byte("init"),
		[]byte("1"),     // # of front-cam
		[]byte("1"),     // # of back-cam
		[]byte("1"),     // # of ALU
		[]byte("1"),     // # of Control Unit
		[]byte("2"),     // # of Register
		[]byte("1"),     // # of Memory
		[]byte("1"),     // # of SSD
		[]byte("1"),     // # of battery
		[]byte("DBS"),   // Bank Account
		[]byte("1000")}) // Bank Balance

	// Check the existence of Entity
	checkEntityUsage(t, stub, "FrontCam0", false)
	checkEntityUsage(t, stub, "BackCam0", false)
	checkEntityUsage(t, stub, "ALU0", false)
	checkEntityUsage(t, stub, "ControlUnit0", false)
	checkEntityUsage(t, stub, "Register0", false)
	checkEntityUsage(t, stub, "Register1", false)
	checkEntityUsage(t, stub, "Memory0", false)
	checkEntityUsage(t, stub, "SSD0", false)
	checkEntityUsage(t, stub, "Battery0", false)

	checkState(t, stub, "DBS", "1000")

	//Make camera
	res := stub.MockInvoke("1", [][]byte{
		[]byte("MakeCamera"), []byte("FrontCam0"),
		[]byte("BackCam0"), []byte("Camera0")})

	if res.Status != shim.OK {
		fmt.Println("Make_Camera failed: ", string(res.Message))
		t.FailNow()
	}
	checkEntityUsage(t, stub, "FrontCam0", true)
	checkEntityUsage(t, stub, "BackCam0", true)
	checkEntityUsage(t, stub, "Camera0", false)

	//Make CPU
	res = stub.MockInvoke("1", [][]byte{
		[]byte("MakeCPU"), []byte("ALU0"),
		[]byte("ControlUnit0"), []byte("Register0"),
		[]byte("Register1"), []byte("CPU0")})

	if res.Status != shim.OK {
		fmt.Println("Make_CPU failed: ", string(res.Message))
		t.FailNow()
	}
	checkEntityUsage(t, stub, "ALU0", true)
	checkEntityUsage(t, stub, "ControlUnit0", true)
	checkEntityUsage(t, stub, "Register0", true)
	checkEntityUsage(t, stub, "Register1", true)
	checkEntityUsage(t, stub, "CPU0", false)

	//make Mainboard
	res = stub.MockInvoke("1", [][]byte{
		[]byte("MakeMainboard"), []byte("CPU0"),
		[]byte("Memory0"), []byte("SSD0"),
		[]byte("Mainboard0")})

	if res.Status != shim.OK {
		fmt.Println("Make_Mainboard failed: ", string(res.Message))
		t.FailNow()
	}
	checkEntityUsage(t, stub, "CPU0", true)
	checkEntityUsage(t, stub, "Memory0", true)
	checkEntityUsage(t, stub, "SSD0", true)
	checkEntityUsage(t, stub, "Mainboard0", false)

	// Assemble IPhone
	res = stub.MockInvoke("1", [][]byte{
		[]byte("Assemble"), []byte("Camera0"),
		[]byte("Battery0"), []byte("Mainboard0"),
		[]byte("IPhone0"), []byte("Manufacturer0")})

	if res.Status != shim.OK {
		fmt.Println("Assemble IPhone failed: ", string(res.Message))
		t.FailNow()
	}
	checkEntityUsage(t, stub, "Camera0", true)
	checkEntityUsage(t, stub, "Battery0", true)
	checkEntityUsage(t, stub, "Mainboard0", true)
	checkIPhoneOwner(t, stub, "IPhone0", "Manufacturer0")

	// Procure IPhone from manufacturer to retailer
	res = stub.MockInvoke("1", [][]byte{
		[]byte("Procure"), []byte("IPhone0"),
		[]byte("Manufacturer0"), []byte("Retailer0")})

	if res.Status != shim.OK {
		fmt.Println("Procure IPhone failed: ", string(res.Message))
		t.FailNow()
	}
	checkIPhoneOwner(t, stub, "IPhone0", "Retailer0")

	// Purchase IPhone using an account
	res = stub.MockInvoke("1", [][]byte{
		[]byte("Purchase"), []byte("IPhone0"),
		[]byte("Customer0"), []byte("DBS"),
		[]byte("Retailer0"), []byte("100")})

	if res.Status != shim.OK {
		fmt.Println("Procure IPhone failed: ", string(res.Message))
		t.FailNow()
	}
	checkIPhoneOwner(t, stub, "IPhone0", "Customer0")
	checkState(t, stub, "DBS", "900")
}
