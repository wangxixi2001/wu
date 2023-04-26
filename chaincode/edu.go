package main

import (
"github.com/hyperledger/fabric-chaincode-go/shim"
"github.com/hyperledger/fabric-protos-go/peer"
"fmt"
"encoding/json"
"bytes"

)


type Certificate struct {
	ObjectType string `json:"docType"`
	AssetName  string `json:"Name"`   
	OwnerID    string `json:"Gender"` 
	Key        string `json:"Key"`
	State      string `json:"Nation"`     
	Version    string `json:"Place"`     
	CertNo     string `json:"CertNo"`     
	Ciphertext string `json:"Graduation"` 
	Note       string `json:"Photo"`     

	Historys []HistoryItem 
}

type HistoryItem struct {
	TxId        string
	Certificate Certificate
}

type CertificateChaincode struct {
}

func (t *CertificateChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

func (t *CertificateChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addCert" {
		return t.addCert(stub, args) 
	} else if fun == "queryCertByCertNoAndName" {
		return t.queryCertByCertNoAndName(stub, args) 
	} else if fun == "queryCertInfoByOwnerID" {
		return t.queryCertInfoByOwnerID(stub, args) 
	} else if fun == "updateCert" {
		return t.updateCert(stub, args) 
	} 
	return shim.Error("指定的函数名称错误")

}

const DOC_TYPE = "eduObj"


func PutCert(stub shim.ChaincodeStubInterface, cert Certificate) ([]byte, bool) {

	cert.ObjectType = DOC_TYPE

	b, err := json.Marshal(cert)
	if err != nil {
		return nil, false
	}


	err = stub.PutState(cert.OwnerID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}



func GetCertInfo(stub shim.ChaincodeStubInterface, ownerID string) (Certificate, bool) {
	var cert Certificate

	b, err := stub.GetState(ownerID)
	if err != nil {
		return cert, false
	}

	if b == nil {
		return cert, false
	}


	err = json.Unmarshal(b, &cert)
	if err != nil {
		return cert, false
	}

	return cert, true
}

func getCertByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}


		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}


func (t *CertificateChaincode) addCert(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var cert Certificate
	err := json.Unmarshal([]byte(args[0]), &cert)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetCertInfo(stub, cert.OwnerID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}

	_, bl := PutCert(stub, cert)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}


func (t *CertificateChaincode) queryCertByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	AssetName := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, AssetName)

	// 查询数据
	result, err := getCertByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及资产名称查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及资产名称没有查询到相关的信息")
	}
	return shim.Success(result)
}


func (t *CertificateChaincode) queryCertInfoByOwnerID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	OwnerID := args[0]
	AssetName := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"OwnerID\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, OwnerID, AssetName)

	// 查询数据
	result, err := getCertByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据AssetName及OwnerID查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的OwnerID及AssetName称没有查询到相关的信息")
	}
	return shim.Success(result)
}


func (t *CertificateChaincode) updateCert(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Certificate
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return shim.Error("反序列化edu信息失败")
	}

	// 根据身份证号码查询信息
	result, bl := GetCertInfo(stub, info.OwnerID)
	if !bl {
		return shim.Error("根据身份证号码查询信息时发生错误")
	}

	result.AssetName = info.AssetName
	result.OwnerID = info.OwnerID
	result.Key = info.Key
	result.State = info.State
	result.Version = info.Version
	result.CertNo = info.CertNo
	result.Ciphertext = info.Ciphertext
	result.Note = info.Note

	_, bl = PutCert(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息更新成功"))
}

func main() {
	err := shim.Start(new(CertificateChaincode))
	if err != nil {
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
