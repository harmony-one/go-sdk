package governance

import "encoding/json"

func submitMessage(address string, content string, sign string) (map[string]interface{}, error) {
	message, err := json.Marshal(map[string]string{
		"address": address,
		"msg":     content,
		"sig":     sign,
	})
	if err != nil {
		return nil, err
	}
	return postAndParse(urlMessage, message)
}
