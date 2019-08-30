package rpc

type reply map[string]interface{}

type T interface {
	SendRPC(RPCMethod, []interface{}) reply
}

type HTTPMessenger struct {
	node string
}

func (M *HTTPMessenger) SendRPC(meth RPCMethod, params []interface{}) reply {
	return RPCRequest(meth, M.node, params)
}

func NewHTTPHandler(node string) *HTTPMessenger {
	// TODO Sanity check the URL for HTTP
	return &HTTPMessenger{node}
}
