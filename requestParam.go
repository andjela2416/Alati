package main

// swagger:parameters deleteConfig
type DeleteConfigRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getConfigById
type GetConfigRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters post createConfig
type RequestConfigBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Config"
	//  required: true
	Body Config `json:"body"`
}

// swagger:parameters post createGroup
type RequestGroupBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/Group"
	//  required: true
	Body Group `json:"body"`
}

// swagger:parameters addConfigToGroup
type AddConfigToGroupRequest struct {
	// Group ID
	// in: path
	GroupId string `json:"g_id"`
	// Config ID
	// in: path
	ConfigId string `json:"c_id"`
}

// swagger:parameters getGroupById
type GetGroupRequest struct {
	// Group ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters deleteGroup
type DeleteGroupRequest struct {
	// Group ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters deleteConfigFromGroup
type DeleteConfigFromGroupRequest struct {
	// Group ID
	// in: path
	GroupId string `json:"g_id"`

	// Config ID
	// in: path
	ConfigId string `json:"c_id"`
}
