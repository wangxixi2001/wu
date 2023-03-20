package main

import (
"github.com/hyperledger/fabric-chaincode-go/shim"
"github.com/hyperledger/fabric-protos-go/peer"
"fmt"
"encoding/json"
"bytes"

)


type Education struct {
	ObjectType	string	`json:"docType"`
	AssetName	string	`json:"Name"`		// 姓名
	OwnerID	string	`json:"Gender"`		// 性别
	State	string	`json:"Nation"`		// 民族
	Version	string	`json:"Place"`		// 籍贯
	CertNo	string	`json:"CertNo"`	// 证书编号
        Ciphertext	string	`json:"Graduation"`	// 毕（结）业
	Note	string	`json:"Photo"`	// 照片

	Historys	[]HistoryItem	// 当前edu的历史记录
}

type HistoryItem struct {
	TxId	string
	Education	Education
}

type EducationChaincode struct {

}

func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addEdu"{
		return t.addEdu(stub, args)		// 添加信息
	}else if fun == "queryEduByCertNoAndName" {
		return t.queryEduByCertNoAndName(stub, args)		// 根据证书编号及姓名查询信息
	}else if fun == "queryEduInfoByEntityID" {
		return t.queryEduInfoByEntityID(stub, args)	// 根据身份证号码及姓名查询详情
	}else if fun == "updateEdu" {
		return t.updateEdu(stub, args)		// 根据证书编号更新信息
	}else if fun == "delEdu"{
		return t.delEdu(stub, args)	// 根据证书编号删除信息
	}

	return shim.Error("指定的函数名称错误")

}


const DOC_TYPE = "eduObj"

// 保存edu
// args: education
func PutEdu(stub shim.ChaincodeStubInterface, edu Education) ([]byte, bool) {

	edu.ObjectType = DOC_TYPE

	b, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}

	// 保存edu状态
	err = stub.PutState(edu.EntityID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据身份证号码查询信息状态
// args: entityID
func GetEduInfo(stub shim.ChaincodeStubInterface, ownerID string) (Education, bool)  {
	var edu Education
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(ownerID)
	if err != nil {
		return edu, false
	}

	if b == nil {
		return edu, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return edu, false
	}

	// 返回结果
	return edu, true
}

// 根据指定的查询字符串实现富查询
func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: educationObject
// 身份证号为 key, Education 为 value
func (t *EducationChaincode) addEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var edu Education
	err := json.Unmarshal([]byte(args[0]), &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetEduInfo(stub, edu.OwnerID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}

	_, bl := PutEdu(stub, edu)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据证书编号及姓名查询信息
// args: CertNo, name
func (t *EducationChaincode) queryEduByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	AssetName := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, AssetName)

	// 查询数据
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及资产名称查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及资产名称没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据身份证号码查询详情（溯源）
// args: entityID
func (t *EducationChaincode) queryEduInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据身份证号码查询edu状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据身份证号码查询信息失败")
	}

	if b == nil {
		return shim.Error("根据身份证号码没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var edu Education
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(edu.OwnerID)
	if err != nil {
		return shim.Error("根据指定的身份证号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisEdu Education
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取edu的历史变更数据失败")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisEdu)

		if hisData.Value == nil {
			var empty Education
			historyItem.Education = empty
		}else {
			historyItem.Education = hisEdu
		}

		historys = append(historys, historyItem)

	}

	edu.Historys = historys

	// 返回
	result, err := json.Marshal(edu)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
}

// 根据身份证号更新信息
// args: educationObject
func (t *EducationChaincode) updateEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Education
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 根据身份证号码查询信息
	result, bl := GetEduInfo(stub, info.OwnerID)
	if !bl{
		return shim.Error("根据身份证号码查询信息时发生错误")
	}

	result.AssetName = info.AssetName
	result.OwnerID= info.OwnerID
	result.State = info.State
	result.Version = info.Version
        result.CertNo = info.CertNo
	result.Ciphertext = info.Ciphertext
	result.Note = info.Note;
	
	_, bl = PutEdu(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据身份证号删除信息（暂不提供）
// args: entityID
func (t *EducationChaincode) delEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息删除成功"))
}

func main(){
	err := shim.Start(new(EducationChaincode))
	if err != nil{
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}

