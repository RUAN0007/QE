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

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SupplyChaincode example simple Chaincode implementation
type TracableChaincode struct {
}

func (cc TracableChaincode) GetLatestWriteTxnForAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Aprovbytes, err := stub.GetState(A + "_prov")
	if err != nil {
		return shim.Error("Failed to get provenance info for " + A)
	}
	if Aprovbytes == nil {
		return shim.Error("Provenance for " + A + " not found")
	}

	var a_prov shim.ProvenanceMeta

	err = json.Unmarshal(Aprovbytes, &a_prov)

	if err != nil {
		return shim.Error("Fail to unmarshal provenance records for " + A)
	}

	return shim.Success([]byte(a_prov.TxID))
}

// SupplyChaincode example simple Chaincode implementation
type SupplyChaincode struct {
	TracableChaincode
}

type Entity struct {
	SerialID string
	Used     bool
}

type Iphone struct {
	SerialID string
	Owner    string
}

func (t *SupplyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	stub.EnableProvenance()

	fmt.Println("Enabling Provenance Tracking")

	_, args := stub.GetFunctionAndParameters()
	var err error

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// Initialize the chaincode
	front_camera_count, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for the number of front cameras")
	}

	// Initing the frontend camera
	for i := 0; i < front_camera_count; i++ {
		var front_camera_serial = "FrontCam" + strconv.Itoa(i)
		var front_camera = Entity{front_camera_serial, false}
		var front_camera_bytes, _ = json.Marshal(front_camera)
		stub.PutState(front_camera_serial, front_camera_bytes)
	}

	back_camera_count, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for the number of back cameras")
	}
	// Initing the backend camera
	for i := 0; i < back_camera_count; i++ {
		var back_camera_serial = "BackCam" + strconv.Itoa(i)
		var back_camera = Entity{back_camera_serial, false}
		var back_camera_bytes, _ = json.Marshal(back_camera)
		stub.PutState(back_camera_serial, back_camera_bytes)
	}

	alu_count, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for the number of alu")
	}

	// Initing the ALU
	for i := 0; i < alu_count; i++ {
		var alu_serial = "ALU" + strconv.Itoa(i)
		var alu = Entity{alu_serial, false}
		var alu_bytes, _ = json.Marshal(alu)
		stub.PutState(alu_serial, alu_bytes)
	}

	control_unit_count, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer for the number of control_unit")
	}

	// Initing the ALU
	for i := 0; i < control_unit_count; i++ {
		var control_unit_serial = "ControlUnit" + strconv.Itoa(i)
		var control_unit = Entity{control_unit_serial, false}
		var control_unit_bytes, _ = json.Marshal(control_unit)
		stub.PutState(control_unit_serial, control_unit_bytes)
	}

	register_count, err := strconv.Atoi(args[4])
	if err != nil {
		return shim.Error("Expecting integer for the number of register")
	}

	// Initing the Register inventory
	for i := 0; i < register_count; i++ {
		var register_serial = "Register" + strconv.Itoa(i)
		var register = Entity{register_serial, false}
		var register_bytes, _ = json.Marshal(register)
		stub.PutState(register_serial, register_bytes)
	}

	memory_count, err := strconv.Atoi(args[5])
	if err != nil {
		return shim.Error("Expecting integer for the number of memory")
	}

	// Initing the memory register
	for i := 0; i < memory_count; i++ {
		var memory_serial = "Memory" + strconv.Itoa(i)
		var memory = Entity{memory_serial, false}
		var memory_bytes, _ = json.Marshal(memory)
		stub.PutState(memory_serial, memory_bytes)
	}

	SSD_count, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error("Expecting integer for the number of SSD")
	}

	// Initing the SSD register
	for i := 0; i < SSD_count; i++ {
		var SSD_serial = "SSD" + strconv.Itoa(i)
		var SSD = Entity{SSD_serial, false}
		var SSD_bytes, _ = json.Marshal(SSD)
		stub.PutState(SSD_serial, SSD_bytes)
	}

	battery_count, err := strconv.Atoi(args[7])
	if err != nil {
		return shim.Error("Expecting integer for the number of battery")
	}

	// Initing the battery register
	for i := 0; i < battery_count; i++ {
		var battery_serial = "Battery" + strconv.Itoa(i)
		var battery = Entity{battery_serial, false}
		var battery_bytes, _ = json.Marshal(battery)
		stub.PutState(battery_serial, battery_bytes)
	}

	// Init the bank account
	bank_account := args[8]
	bank_balance, err := strconv.Atoi(args[9])
	if err != nil {
		return shim.Error("Expecting integer for bank balance. ")
	}

	stub.PutState(bank_account, []byte(strconv.Itoa(bank_balance)))

	return shim.Success(nil)
}

func (t *SupplyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	stub.EnableProvenance()
	function, args := stub.GetFunctionAndParameters()
	if function == "lastWrtTxn" {
		// Make payment of X units from A to B
		return t.GetLatestWriteTxnForAsset(stub, args)
	} else if function == "MakeCamera" {
		return t.make_camera(stub, args)
	} else if function == "MakeCPU" {
		return t.make_cpu(stub, args)
	} else if function == "MakeMainboard" {
		return t.make_mainboard(stub, args)
	} else if function == "Assemble" {
		return t.assemble_iphone(stub, args)
	} else if function == "Procure" {
		return t.procure(stub, args)
	} else if function == "Purchase" {
		return t.purchase(stub, args)
	} else if function == "Query" {
		return t.query(stub, args)
	}
	return shim.Error("Invalid invoke function name.")
}

func (t *SupplyChaincode) purchase(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	iphone_serial := args[0]
	customer := args[1]
	bank_account := args[2]
	retailer := args[3]
	price, err := strconv.Atoi(args[4])

	if err != nil {
		return shim.Error("Expecting integer value for price ")
	}

	iphone_bytes, err := stub.GetState(iphone_serial)
	if err != nil {
		return shim.Error("No Manufactured iphone with ID " + iphone_serial)
	}
	var iphone Iphone
	err = json.Unmarshal(iphone_bytes, &iphone)
	if err != nil {
		return shim.Error("Cannot unmarshal iPhone with ID " + iphone_serial)
	}

	if iphone.Owner != retailer {
		return shim.Error("Iphone with ID " + iphone_serial + " is not owned by retailer " + retailer)
	}

	bank_balance_raw, err := stub.GetState(bank_account)
	if err != nil {
		return shim.Error("Cannot find account " + bank_account)
	}

	bank_balance, err := strconv.Atoi(string(bank_balance_raw))
	if err != nil {
		return shim.Error("Expect integer for bank balance " + bank_account)
	}

	if bank_balance < price {
		return shim.Error("The account does not have enough balance. Purchasing fails")
	}

	bank_balance -= price
	// Put back the bank balance
	err = stub.PutState(bank_account, []byte(strconv.Itoa(bank_balance)))
	if err != nil {
		return shim.Error(err.Error())
	}

	iphone.Owner = customer
	iphone_bytes, err = json.Marshal(iphone)
	err = stub.PutState(iphone_serial, iphone_bytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SupplyChaincode) procure(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	iphone_serial := args[0]
	manufactuerer := args[1]
	retailer := args[2]

	iphone_bytes, err := stub.GetState(iphone_serial)
	if err != nil {
		return shim.Error("No Manufactured iphone with ID " + iphone_serial)
	}
	var iphone Iphone
	err = json.Unmarshal(iphone_bytes, &iphone)
	if err != nil {
		return shim.Error("Cannot unmarshal iPhone with ID " + iphone_serial)
	}

	if iphone.Owner != manufactuerer {
		return shim.Error("Iphone with ID " + iphone_serial + " is not owned by manufacturer " + manufactuerer)
	}

	iphone.Owner = retailer
	iphone_bytes, _ = json.Marshal(iphone)
	stub.PutState(iphone_serial, iphone_bytes)

	return shim.Success(nil)
}

func (t *SupplyChaincode) assemble_iphone(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	camera_serial := args[0]

	// Retrieve the camera
	camera_bytes, err := stub.GetState(camera_serial)
	if err != nil {
		return shim.Error("No camera with ID " + camera_serial)
	}
	var camera Entity
	err = json.Unmarshal(camera_bytes, &camera)
	if err != nil {
		return shim.Error("Cannot unmarshal camera with ID " + camera_serial)
	}
	if camera.Used {
		return shim.Error("Camera with ID " + camera_serial + " is used. ")
	}
	camera.Used = true
	camera_bytes, _ = json.Marshal(camera)
	stub.PutState(camera_serial, camera_bytes)

	// Retrive the battery
	battery_serial := args[1]
	battery_bytes, err := stub.GetState(battery_serial)
	if err != nil {
		return shim.Error("No battery with ID " + battery_serial)
	}
	var battery Entity
	err = json.Unmarshal(battery_bytes, &battery)
	if err != nil {
		return shim.Error("Cannot unmarshal battery with ID " + battery_serial)
	}
	if battery.Used {
		return shim.Error("Battery with ID " + battery_serial + " is used. ")
	}
	battery.Used = true
	battery_bytes, _ = json.Marshal(battery)
	stub.PutState(battery_serial, battery_bytes)

	// Retrive the mainboard
	mainboard_serial := args[2]
	mainboard_bytes, err := stub.GetState(mainboard_serial)
	if err != nil {
		return shim.Error("No mainboard with ID " + mainboard_serial)
	}
	var mainboard Entity
	err = json.Unmarshal(mainboard_bytes, &mainboard)
	if err != nil {
		return shim.Error("Cannot unmarshal Mainboard with ID " + mainboard_serial)
	}
	if mainboard.Used {
		return shim.Error("Mainboard with ID " + mainboard_serial + " is used. ")
	}
	mainboard.Used = true
	mainboard_bytes, _ = json.Marshal(mainboard)
	stub.PutState(mainboard_serial, mainboard_bytes)

	// Put the manufactured mainboard
	iphone_serial := args[3]
	manufacturer := args[4]
	iphone := Iphone{iphone_serial, manufacturer}
	iphone_bytes, _ := json.Marshal(iphone)
	stub.PutState(iphone_serial, iphone_bytes)

	return shim.Success(nil)
}

func (t *SupplyChaincode) make_camera(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	front_cam_serial := args[0]

	// Retrive the front camera asset
	front_cam_bytes, err := stub.GetState(front_cam_serial)
	if err != nil {
		return shim.Error("No front camera with ID " + front_cam_serial)
	}
	var front_cam Entity
	err = json.Unmarshal(front_cam_bytes, &front_cam)
	if err != nil {
		return shim.Error("Cannot unmarshal front camera with ID " + front_cam_serial)
	}
	if front_cam.Used {
		return shim.Error("Front Camera with ID is ")
	}
	front_cam.Used = true
	front_cam_bytes, _ = json.Marshal(front_cam)
	stub.PutState(front_cam_serial, front_cam_bytes)

	// Retrive the back camera asset
	back_cam_serial := args[1]
	back_cam_bytes, err := stub.GetState(back_cam_serial)
	if err != nil {
		return shim.Error("No back camera with ID " + back_cam_serial)
	}
	var back_cam Entity
	err = json.Unmarshal(back_cam_bytes, &back_cam)
	if err != nil {
		return shim.Error("Cannot unmarshal back camera with ID " + back_cam_serial)
	}
	back_cam.Used = true
	back_cam_bytes, _ = json.Marshal(back_cam)
	stub.PutState(back_cam_serial, back_cam_bytes)

	// Put the manufactured camera
	camera_serial := args[2]
	var camera = Entity{camera_serial, false}
	camera_bytes, _ := json.Marshal(camera)
	stub.PutState(camera_serial, camera_bytes)

	return shim.Success(nil)
}

func (t *SupplyChaincode) make_cpu(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	alu_serial := args[0]

	// Retrive the alu asset
	alu_bytes, err := stub.GetState(alu_serial)
	if err != nil {
		return shim.Error("No ALU with ID " + alu_serial)
	}
	var alu Entity
	err = json.Unmarshal(alu_bytes, &alu)
	if err != nil {
		return shim.Error("Cannot unmarshal ALU with ID " + alu_serial)
	}
	if alu.Used {
		return shim.Error("ALU with ID " + alu_serial + " is used. ")
	}

	alu.Used = true
	alu_bytes, _ = json.Marshal(alu)
	stub.PutState(alu_serial, alu_bytes)

	// Retrive the control unit asset
	control_unit_serial := args[1]
	control_unit_bytes, err := stub.GetState(control_unit_serial)
	if err != nil {
		return shim.Error("No control unit with ID " + control_unit_serial)
	}
	var control_unit Entity
	err = json.Unmarshal(control_unit_bytes, &control_unit)
	if err != nil {
		return shim.Error("Cannot unmarshal control unit with ID " + control_unit_serial)
	}
	if control_unit.Used {
		return shim.Error("Control Unit with ID " + control_unit_serial + " is used. ")
	}
	control_unit.Used = true
	control_unit_bytes, _ = json.Marshal(control_unit)
	stub.PutState(control_unit_serial, control_unit_bytes)

	// Retrive the register1 asset
	register1_serial := args[2]
	register1_bytes, err := stub.GetState(register1_serial)
	if err != nil {
		return shim.Error("No register with ID " + register1_serial)
	}
	var register1 Entity
	err = json.Unmarshal(register1_bytes, &register1)
	if err != nil {
		return shim.Error("Cannot unmarshal register with ID " + register1_serial)
	}
	if register1.Used {
		return shim.Error("Register with ID " + register1_serial + " is used. ")
	}
	register1.Used = true
	register1_bytes, _ = json.Marshal(register1)
	stub.PutState(register1_serial, register1_bytes)

	// Retrive the register2 asset
	register2_serial := args[3]
	register2_bytes, err := stub.GetState(register2_serial)
	if err != nil {
		return shim.Error("No register with ID " + register2_serial)
	}
	var register2 Entity
	err = json.Unmarshal(register2_bytes, &register2)
	if err != nil {
		return shim.Error("Cannot unmarshal register with ID " + register2_serial)
	}
	if register2.Used {
		return shim.Error("Register with ID " + register2_serial + " is used. ")
	}
	register2.Used = true
	register2_bytes, _ = json.Marshal(register2)
	stub.PutState(register2_serial, register2_bytes)

	// Put the manufactured cpu
	cpu_serial := args[4]
	var cpu = Entity{cpu_serial, false}
	cpu_bytes, _ := json.Marshal(cpu)
	stub.PutState(cpu_serial, cpu_bytes)
	return shim.Success(nil)
}

func (t *SupplyChaincode) make_mainboard(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	cpu_serial := args[0]

	// Retrive the cpu
	cpu_bytes, err := stub.GetState(cpu_serial)
	if err != nil {
		return shim.Error("No CPU with ID " + cpu_serial)
	}
	var cpu Entity
	err = json.Unmarshal(cpu_bytes, &cpu)
	if err != nil {
		return shim.Error("Cannot unmarshal CPU with ID " + cpu_serial)
	}
	if cpu.Used {
		return shim.Error("CPU with ID " + cpu_serial + " is used. ")
	}
	cpu.Used = true
	cpu_bytes, _ = json.Marshal(cpu)
	stub.PutState(cpu_serial, cpu_bytes)

	// Retrive the memory
	memory_serial := args[1]
	memory_bytes, err := stub.GetState(memory_serial)
	if err != nil {
		return shim.Error("No memory with ID " + memory_serial)
	}
	var memory Entity
	err = json.Unmarshal(memory_bytes, &memory)
	if err != nil {
		return shim.Error("Cannot unmarshal memory with ID " + memory_serial)
	}
	if memory.Used {
		return shim.Error("Memory with ID " + cpu_serial + " is used. ")
	}
	memory.Used = true
	memory_bytes, _ = json.Marshal(memory)
	stub.PutState(memory_serial, memory_bytes)

	// Retrive the SSD
	SSD_serial := args[2]
	SSD_bytes, err := stub.GetState(SSD_serial)
	if err != nil {
		return shim.Error("No SSD with ID " + SSD_serial)
	}
	var SSD Entity
	err = json.Unmarshal(SSD_bytes, &SSD)
	if err != nil {
		return shim.Error("Cannot unmarshal SSD with ID " + SSD_serial)
	}
	if SSD.Used {
		return shim.Error("SSD with ID " + cpu_serial + " is used. ")
	}
	SSD.Used = true
	SSD_bytes, _ = json.Marshal(SSD)
	stub.PutState(SSD_serial, SSD_bytes)

	// Put the manufactured mainboard
	mainboard_serial := args[3]
	var mainboard = Entity{mainboard_serial, false}
	mainboard_bytes, _ := json.Marshal(mainboard)
	stub.PutState(mainboard_serial, mainboard_bytes)

	return shim.Success(nil)
}

func (t *SupplyChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 argument")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SupplyChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
