package main

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

type VMPoint struct {
    UserId    int    `json:"user_id"`
    CorpName  string `json:"corporate_name"`
    Point     int    `json:"point"`
}

type SimpleChaincode struct {
}

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, functio string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }

    return nil, nil
}

func (t *SimpleChaincode) init_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var err error
    id := args[0]
    user_id, _ := strconv.Atoi(args[1])
    corporate_name := strings.ToLower(args[2])
    point, _ := strconv.Atoi(args[3])

    str := `{"user_id": ` + strconv.Itoa(user_id) + `,"corporate_name": "` + corporate_name + `","point": ` + strconv.Itoa(point) + `}`

    fmt.Println(str)

    err = stub.PutState(id, []byte(str))
    if err != nil {
        return nil, err
    }
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("invoke is running" + function)

    if function == "init" {
        return t.Init(stub, "Init", args)
    } else if function == "add" {
        return t.init_user(stub, args)
    } else if function == "deposit" {
        point, _ := strconv.Atoi(args[1])
        VMPointAsByte, err := stub.GetState(args[0])
        if err != nil {
            return nil, err
        }

        res := VMPoint{}
        json.Unmarshal(VMPointAsByte, &res)
        res.Point = res.Point + point

        jsonAsBytes, _ := json.Marshal(res)        
        stub.PutState(args[0], jsonAsBytes)

        return nil, nil
    
    } else if function == "transfer" {
        VMPointAsByte_dest, err := stub.GetState(args[0])
        if err != nil {
            return nil, err
        }
        VMPointAsByte_src, err := stub.GetState(args[1])
        if err != nil {
            return nil, err
        }

        res_dest := VMPoint{}
        json.Unmarshal(VMPointAsByte_dest, &res_dest)
        res_src := VMPoint{}
        json.Unmarshal(VMPointAsByte_src, &res_src)

        point_diff, _ := strconv.Atoi(args[2])

        res_dest.Point = res_dest.Point - point_diff
        res_src.Point = res_src.Point + point_diff

        jsonAsBytes_dest, _ := json.Marshal(res_dest)
        jsonAsBytes_src, _ := json.Marshal(res_src)

        stub.PutState(args[0], jsonAsBytes_dest)
        stub.PutState(args[1], jsonAsBytes_src)

        return nil, nil

    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("query is running" + function)
    if function == "read" {
        return t.read(stub, args)
    }
    fmt.Println("query did not find func: " + function)

    return nil, errors.New("Received unknown function query:" + function)
} 
