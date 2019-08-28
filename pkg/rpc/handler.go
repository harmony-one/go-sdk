package rpc

// "github.com/pkg/errors"

// Okay so the interface is just methods
type T interface {
	SendRPC(RPCMethod, interface{}) (string, error)
}

// And the struct is data
type HTTPMessenger struct {
	error
	nodePath string
	queryID  int
}

var (
	HTTPHandler T
)

func init() {
	HTTPHandler = HTTPMessenger{nodePath: "http://localhost:9500"}
}

func (handler HTTPMessenger) SendRPC(method RPCMethod, params interface{}) (string, error) {
	return "!23", nil
}
