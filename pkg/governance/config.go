package governance

type governanceApi string

const backendAddress = "https://snapshot.hmny.io/api/"

const (
	_ governanceApi = ""

	urlListSpace                           = backendAddress + "spaces"
	urlListProposalsBySpace                = backendAddress + "%s/proposals"
	urlListProposalsVoteBySpaceAndProposal = backendAddress + "%s/proposal/%s"
	urlMessage                             = backendAddress + "message"
	urlGetValidatorsInTestNet              = "https://api.stake.hmny.io/networks/testnet/validators"
	urlGetValidatorsInMainNet              = "https://api.stake.hmny.io/networks/mainnet/validators"
	urlGetProposalInfo                     = "https://gateway.ipfs.io/ipfs/%s"
)
