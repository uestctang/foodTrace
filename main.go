
package main

import (
	"2020-Winter-vacation-homework/trustPlatform/user"
	"foodTrace/trace"
	"fmt"


	"fabric/core/chaincode/shim"
pb "fabric/protos/peer"

"strings"
)

// 主模块
type TrustPlatformCC struct {
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(trace.Food))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ===================================================================================
// 初始化函数，如果有需要初始化的模块请在自己模块中定义init函数，在此调用，不要直接写在主模块中
// ===================================================================================
func (t *TrustPlatformCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("TrustPlatformCC Init")
	// example
	trace.Food.Init(stub)
	return shim.Success([]byte("Init success"))
}

// ===================================================================================
// 调用合约入口，按模块划分，例如我们有用户、组织、数据分享三大模块（可以添加），按照路径层级传参
// 例如 用户声明新属性 /user/attribute/add args: userName, attributePublicKey
// user为用户模块，attribute为要操作的用户属性，add代表新建/声明，两个args按指定顺序传入
// ===================================================================================
func (t *TrustPlatformCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("\nTrustPlatformCC Invoke")
	function, _ := stub.GetFunctionAndParameters()
	if strings.HasPrefix(function, "/user") {
		return user.Invoke(stub)
	}

	return shim.Error("Invalid invoke function name. Expecting \"/user\" ")
}
