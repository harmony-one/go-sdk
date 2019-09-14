package rpc

type Reply map[string]interface{}

type T interface {
	SendRPC(RPCMethod, []interface{}) (Reply, error)
}

type HTTPMessenger struct {
	node string
}

func (M *HTTPMessenger) SendRPC(meth RPCMethod, params []interface{}) (Reply, error) {
	return Request(meth, M.node, params)
}

func NewHTTPHandler(node string) *HTTPMessenger {
	// TODO Sanity check the URL for HTTP
	return &HTTPMessenger{node}
}
