package task

const ProtocolV1 = "1"
const ProtocolV1_1 = "1.1"
const ProtocolV2 = "2"

const HydroStartBlockNumberV1 = 6885289
const HydroExchangeAddressV1 = "0x2cB4B49C0d6E9db2164d94Ce48853BF77C4D883E"
const HydroMatchTopicV1 = "0xdcc6682c66bde605a9e21caeb0cb8f1f6fbd5bbfb2250c3b8d1f43bb9b06df3f"
const HydroExchangeABIV1 = `[
	{
	  "anonymous": false,
	  "inputs": [
		{ "indexed": false, "name": "baseToken", "type": "address" },
		{ "indexed": false, "name": "quoteToken", "type": "address" },
		{ "indexed": false, "name": "relayer", "type": "address" },
		{ "indexed": false, "name": "maker", "type": "address" },
		{ "indexed": false, "name": "taker", "type": "address" },
		{ "indexed": false, "name": "baseTokenAmount", "type": "uint256" },
		{ "indexed": false, "name": "quoteTokenAmount", "type": "uint256" },
		{ "indexed": false, "name": "makerFee", "type": "uint256" },
		{ "indexed": false, "name": "takerFee", "type": "uint256" },
		{ "indexed": false, "name": "makerGasFee", "type": "uint256" },
		{ "indexed": false, "name": "makerRebate", "type": "uint256" },
		{ "indexed": false, "name": "takerGasFee", "type": "uint256" }
	  ],
	  "name": "Match",
	  "type": "event"
	}
  ]`

const HydroStartBlockNumberV1_1 = 7454912
const HydroExchangeAddressV1_1 = "0xE2a0BFe759e2A4444442Da5064ec549616FFF101"
const HydroMatchTopicV1_1 = "0xd3ac06c3b34b93617ba2070b8b7a925029035b3f30fecd2d0fa8e5845724f310"
const HydroExchangeABIV1_1 = `[
  {
    "anonymous": false,
    "inputs": [
      {
        "components": [
          { "name": "baseToken", "type": "address" },
          { "name": "quoteToken", "type": "address" },
          { "name": "relayer", "type": "address" }
        ],
        "indexed": false,
        "name": "addressSet",
        "type": "tuple"
      },
      {
        "components": [
          { "name": "maker", "type": "address" },
          { "name": "taker", "type": "address" },
          { "name": "buyer", "type": "address" },
          { "name": "makerFee", "type": "uint256" },
          { "name": "makerRebate", "type": "uint256" },
          { "name": "takerFee", "type": "uint256" },
          { "name": "makerGasFee", "type": "uint256" },
          { "name": "takerGasFee", "type": "uint256" },
          { "name": "baseTokenFilledAmount", "type": "uint256" },
          { "name": "quoteTokenFilledAmount", "type": "uint256" }
        ],
        "indexed": false,
        "name": "result",
        "type": "tuple"
      }
    ],
    "name": "Match",
    "type": "event"
  }
]`

const HydroStartBlockNumberV2 = 8399662
const HydroExchangeAddressV2 = "0x241e82C79452F51fbfc89Fac6d912e021dB1a3B7"
const HydroMatchTopicV2 = "0x6bf96fcc2cec9e08b082506ebbc10114578a497ff1ea436628ba8996b750677c"

// const HydroExchangeABIV2 = ``
