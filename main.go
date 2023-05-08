package main

import (
	"encoding/json"
	"fmt"
	"wu/sdkInit"
	"wu/service"
	"wu/web"
	"wu/web/controller"
	"os"

)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/wu/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/wu/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH") + "/src/wu/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> Setting the chain code status via the chain code external service ......")

	cert := service.Certificate{
		AssetName:  "Labor Contract01",
		OwnerID:    "101",
		Key:        "abc&1*~#^2^#s0^=)^^7%b34",
		State:      "valid",
		Version:    "1.0",
		CertNo:     "111",
		Ciphertext: "uQyiNNQr5tZMuNpZqooMkg==",
		Note:       ".......",
	}

	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err != nil {
		fmt.Println()
		os.Exit(-1)
	}
	msg, err := serviceSetup.SaveCert(cert)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("The message was successfully posted with the transaction number: " + msg)
	}

	result, err := serviceSetup.FindCertByCertNoAndName("111","Labor Contract01")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var cert service.Certificate
		json.Unmarshal(result, &cert)
		fmt.Println("OwnerID query info success：")
		fmt.Println(cert)
	}

	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}
