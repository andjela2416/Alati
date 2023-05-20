package main

// swagger:parameters deletePost
type DeleteConfigRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters getPostById
type GetConfigRequest struct {
	// Post ID
	// in: path
	Id string `json:"id"`
}

// swagger:parameters post createPost
type RequestConfigBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestPost"
	//  required: true
	Body configServer `json:"body"`
}

// swagger:parameters addConfigToGroup
type AddConfigToGroupRequest struct {
	// Group ID
	// Config ID
	// in: path
	GroupId  string `json:"group-id"`
	ConfigId string `json:"config-id"`
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
	// Config ID
	// in: path
	GroupId  string `json:"group-id"`
	ConfigId string `json:"config-id"`
}
