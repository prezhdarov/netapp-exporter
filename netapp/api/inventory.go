package ontapi

import "encoding/json"

func (na *ONTAP) listNodes(loginData map[string]interface{}) (map[string]string, error) {

	type nodeInventory struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}

	type Response struct {
		Records []nodeInventory `json:"records"`
	}

	var data *Response

	nodes := make(map[string]string)

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = "/api/cluster/nodes"
	body, err := na.Get(loginData, extraConfig)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*body.(*[]byte), &data)
	if err != nil {
		return nil, err
	}

	for _, node := range data.Records {
		nodes[node.UUID] = node.Name
	}

	return nodes, nil
}
