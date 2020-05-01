package trace

import (
	"encoding/json"
	"fabric/core/chaincode/shim"
	pb "fabric/protos/peer"
	"fmt"
	"foodTrace/model"
)

type Food struct {
}

func (a *Food) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Food Init")
	return shim.Success([]byte("Init success!"))
}

func (a *Food) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()

	switch fn {
	case "addProInfo":
		return a.addProInfo(stub, args)
	case "addIngInfo":
		return a.addIngInfo(stub, args)
	case "getFoodInfo":
		return a.getFoodInfo(stub, args)
	case "addLogInfo":
		return a.addLogInfo(stub, args)
	case "getProInfo":
		return a.getProInfo(stub, args)
	case "getLogInfo":
		return a.getLogInfo(stub, args)
	case "getIngInfo":
		return a.getIngInfo(stub, args)
	case "getLogInfo_l":
		return a.getLogInfo_l(stub, args)
	}

	return shim.Error(fmt.Sprintf("unsupported function: %s", fn))
}

func (a *Food) addProInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var FoodInfos model.FoodInfo

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodInfos.FoodID = args[0]
	if FoodInfos.FoodID == "" {
		return shim.Error("FoodID can not be empty.")
	}

	FoodInfos.FoodProInfo.FoodName = args[1]
	FoodInfos.FoodProInfo.FoodSpec = args[2]
	FoodInfos.FoodProInfo.FoodMFGDate = args[3]
	FoodInfos.FoodProInfo.FoodEXPDate = args[4]
	FoodInfos.FoodProInfo.FoodLOT = args[5]
	FoodInfos.FoodProInfo.FoodQSID = args[6]
	FoodInfos.FoodProInfo.FoodMFRSName = args[7]
	FoodInfos.FoodProInfo.FoodProPrice = args[8]
	FoodInfos.FoodProInfo.FoodProPlace = args[9]
	ProInfosJSONasBytes, err := json.Marshal(FoodInfos)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(FoodInfos.FoodID, ProInfosJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (a *Food) addIngInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var FoodInfos model.FoodInfo
	var IngInfoitem model.IngInfo

	if (len(args)-1)%2 != 0 || len(args) == 1 {
		return shim.Error("Incorrect number of arguments")
	}

	FoodID := args[0]
	for i := 1; i < len(args); {
		IngInfoitem.IngID = args[i]
		IngInfoitem.IngName = args[i+1]
		FoodInfos.FoodIngInfo = append(FoodInfos.FoodIngInfo, IngInfoitem)
		i = i + 2
	}

	FoodInfos.FoodID = FoodID
	/*  FoodInfos.FoodIngInfo = foodIngInfo*/
	IngInfoJsonAsBytes, err := json.Marshal(FoodInfos)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(FoodInfos.FoodID, IngInfoJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func (a *Food) addLogInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var FoodInfos model.FoodInfo

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodInfos.FoodID = args[0]
	if FoodInfos.FoodID == "" {
		return shim.Error("FoodID can not be empty.")
	}
	FoodInfos.FoodLogInfo.LogDepartureTm = args[1]
	FoodInfos.FoodLogInfo.LogArrivalTm = args[2]
	FoodInfos.FoodLogInfo.LogMission = args[3]
	FoodInfos.FoodLogInfo.LogDeparturePl = args[4]
	FoodInfos.FoodLogInfo.LogDest = args[5]
	FoodInfos.FoodLogInfo.LogToSeller = args[6]
	FoodInfos.FoodLogInfo.LogStorageTm = args[7]
	FoodInfos.FoodLogInfo.LogMOT = args[8]
	FoodInfos.FoodLogInfo.LogCopName = args[9]
	FoodInfos.FoodLogInfo.LogCost = args[10]

	LogInfosJSONasBytes, err := json.Marshal(FoodInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(FoodInfos.FoodID, LogInfosJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (a *Food) getFoodInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(FoodID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var foodAllinfo model.FoodAllInfo

	for resultsIterator.HasNext() {
		var FoodInfos model.FoodInfo
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(response.Value, &FoodInfos)
		if FoodInfos.FoodProInfo.FoodName != "" {
			foodAllinfo.FoodProInfo = FoodInfos.FoodProInfo
		} else if FoodInfos.FoodIngInfo != nil {
			foodAllinfo.FoodIngInfo = FoodInfos.FoodIngInfo
		} else if FoodInfos.FoodLogInfo.LogMission != "" {
			foodAllinfo.FoodLogInfo = append(foodAllinfo.FoodLogInfo, FoodInfos.FoodLogInfo)
		}

	}

	jsonsAsBytes, err := json.Marshal(foodAllinfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsonsAsBytes)
}

func (a *Food) getProInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(FoodID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var foodProInfo model.ProInfo

	for resultsIterator.HasNext() {
		var FoodInfos model.FoodInfo
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(response.Value, &FoodInfos)
		if FoodInfos.FoodProInfo.FoodName != "" {
			foodProInfo = FoodInfos.FoodProInfo
			continue
		}
	}
	jsonsAsBytes, err := json.Marshal(foodProInfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

func (a *Food) getIngInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(FoodID)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var foodIngInfo []model.IngInfo
	for resultsIterator.HasNext() {
		var FoodInfos model.FoodInfo
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(response.Value, &FoodInfos)
		if FoodInfos.FoodIngInfo != nil {
			foodIngInfo = FoodInfos.FoodIngInfo
			continue
		}
	}
	jsonsAsBytes, err := json.Marshal(foodIngInfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

func (a *Food) getLogInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var LogInfos []model.LogInfo

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	FoodID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(FoodID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		var FoodInfos model.FoodInfo
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(response.Value, &FoodInfos)
		if FoodInfos.FoodLogInfo.LogMission != "" {
			LogInfos = append(LogInfos, FoodInfos.FoodLogInfo)
		}
	}
	jsonsAsBytes, err := json.Marshal(LogInfos)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

func getLogInfo_l(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Loginfo model.LogInfo

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	FoodID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(FoodID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		var FoodInfos model.FoodInfo
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(response.Value, &FoodInfos)
		if FoodInfos.FoodLogInfo.LogMission != "" {
			Loginfo = FoodInfos.FoodLogInfo
			continue
		}
	}
	jsonsAsBytes, err := json.Marshal(Loginfo)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonsAsBytes)
}

func main() {
	err := shim.Start(new(Food))
	if err != nil {
		fmt.Printf("Error starting Food chaincode: %s ", err)
	}
}
