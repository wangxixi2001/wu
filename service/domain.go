
package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"wu/sdkInit"
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

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("Registration Chain Code Event Failure: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Receiving a chain code event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("The corresponding chain code event cannot be received based on the specified event ID(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}
