package algosdk

import (
	"context"
	"crypto/ed25519"
	json "encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

const algodAddress = ""
const psToken = ""//Purestake API key

// Networks data structure to hold name, server, label, explorer


type Networks struct {
	name  string
	server url
	label string
	explorer url
}

/*testnet := Networks{'testnet', 'https://testnet-algorand.api.purestake.io/ps1',
						'TESTNET', 'https://goalseeker.purestake.io/algorand/testnet'}
						
mainnet := Networks{'mainnet', 'https://mainnet-algorand.api.purestake.io/ps1',
						'MAINNET', 'https://goalseeker.purestake.io/algorand/mainnet'}
*/
type AlgoSDK struct {
	client
	network
	Networks
}

func (a *AlgoSDK) setNetwork(name) {

	// let's initialize testnet and mainnet networks
	testnet := Networks{'testnet', 'https://testnet-algorand.api.purestake.io/ps1',
						'TESTNET', 'https://goalseeker.purestake.io/algorand/testnet'}
						
	mainnet := Networks{'mainnet', 'https://mainnet-algorand.api.purestake.io/ps1',
						'MAINNET', 'https://goalseeker.purestake.io/algorand/mainnet'}
	
	if name == testnet.name {
		a.network := testnet
	} else if name == mainnet.name {
		a.network := mainnet.server 
	}
	return a.network
}

func (a *AlgoSDK) getExplorer() {
	const algodAddress := setNetWork()

}

func (a *AlgoSDK) getCurrentNetwork() {
	return a.network
}

func (a *AlgoSDK) selectNetwork(name) {
	for network, values := range getNetworks() {
		if network.name === name {
			setAlgodClient(network.server)
		}
	}
}

func (n *Networks) getNetworks() {
	testNet := Networks{'testnet', 'https://testnet-algorand.api.purestake.io/ps1',
						'TESTNET', 'https://goalseeker.purestake.io/algorand/testnet'}
						
	mainNet := Networks{'mainnet', 'https://mainnet-algorand.api.purestake.io/ps1',
						'MAINNET', 'https://goalseeker.purestake.io/algorand/mainnet'}

	return testNet, mainNet
}

func (a *AlgoSDK) setAlgodClient(Networks.server) {
	// Set up Headers first
	testNet, mainNet := getNetworks()
	const algodAddressTestNet := testNet.server
	const algodAddressMainNet := MainNet.server
	var headers []*algod.Header
	headers = append(headers, &algod.Header{"X-API-Key", psToken})
	algodClient, err := algod.MakeClientWithHeaders(algodAddressTestNet, "", headers )
	if err != nil {
			fmt.Printf("failed to make algod client: %s\n", err)
			return
	}
	client := algodClient
	return client
}

func (a *AlgoSDK) getCurrentNetwork() {
	return a.network
}

/* Convert from Javascript: 
selectNetwork(name) {
        const networks = this.getNetworks();
        networks.forEach((network) => {
            if (network.name === name) {
                this.network = network;
                this.setClient(network.server);
            }
        })
    }
	*/
func (a AlgoSDK) selectNetwork(networks) {
	//selectNetwork(networks == networks[testnet])
}


/* Convert from Javascipt:
getExplorer() {
        const network = this.getCurrentNetwork();
        return this.network.explorer;
    }
*/
func(a *AlgoSDK) getExplorer() {
	const network = getCurrentNetwork();
	
	return network.explorer;
}

func (a *AlgoSDK) getAssetUrl(id) {
	return getExplorer() + '/asset/' + id;
}

func (a *AlgoSDK) getCurrentNetwork() {
	return network
}

func (a *AlgoSDK) setClient(server) {
	p := *AlgoSDK
	//p.client = new algod.Client(token, server, port);
}

func (a *AlgoSDK) getClient() {
	return //p.client
}

/* Implement this as well:
mnemonicToSecretKey(mnemonic) {
        return sdk.mnemonicToSecretKey(mnemonic);
    }
*/

func (a AlgoSDK) getAccountInformation(address) {
	return getClient().accountInformation(address)
}

func (a AlgoSDK) getAssetInformation(assetID) {
	return getClient().assetInformation(assetID);
}

func (a AlgoSDK) getChangingParams() {
	const cp = {
		fee: 0,
		firstRound: 0,
		lastRound: 0,
		genID: "",
		genHash: ""
	}

	let params = // await this.getClient().getTransactionParams();
	// cp.firstRound = params.lastRound;
	// cp.lastRound = cp.firstRound + parseInt(1000);
	// let sFee = await this.getClient().suggestedFee();
	/* cp.fee = sFee.fee;
	cp.genID = params.genesisID;
	cp.genHash = params.genesishashb64;

	return cp; */
}

func(a AlgoSDK) waitForConfirmation(txId) {
	/* let lastRound = (await this.getClient().status()).lastRound;
	while (true) {
		const pendingInfo = await this.getClient().pendingTransactionInformation(txId);
		if (pendingInfo.round !== null && pendingInfo.rround > 0) {
			// The completedt transaction
			console.log("Transaction " + pendingInfo.tx + " confirmed in round " + pendingInfo.round);
			break;
		}
		lastRound++;
		await this.getClient().statusAfterBlock(lastRound);
		*/
	}
}


