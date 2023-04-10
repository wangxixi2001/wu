package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"encoding/json"
	"fmt"
)

func (t *ServiceSetup) SaveCert(cert Certificate) (string, error) {

	eventID := "eventAddEdu"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(cert)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addCert", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) FindCertInfoByOwnerID(ownerID, assetName string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryCertInfoByOwnerID", Args: [][]byte{[]byte(ownerID), []byte(assetName)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) FindCertByCertNoAndName(certNo, name string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "queryCertByCertNoAndName", Args: [][]byte{[]byte(certNo), []byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

func (t *ServiceSetup) ModifyCert(cert Certificate) (string, error) {

	eventID := "eventModifyEdu"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将edu对象序列化成为字节数组
	b, err := json.Marshal(cert)
	if err != nil {
		return "", fmt.Errorf("Error while serialising the specified cert object")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "updateCert", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) DelCert(ownerID string) (string, error) {

	eventID := "eventDelEdu"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "delCert", Args: [][]byte{[]byte(ownerID), []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}
