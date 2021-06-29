package purestake

import (
	"context"
	"fmt"
	
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
)

const algodAddress = "https://testnet-algorand.api.purestake.io/ps2"
const psToken = ""

const port = '';

struct psAlgoSDK {
}

// Data Structure to store the two networks, the
// data they send and how to access them
// This data structure will be instantiated in
// the psAlgoSDK struct as well

// Network parameters are:
// name, server, label, and explorer

networks := map[string]URL{
	name: 'testnet',
	server: 'https://testnet-algorand.api.purestake'
}


// method to instantiate psAlgoSDK
// take network name as parameter

