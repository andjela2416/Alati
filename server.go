package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type configServer struct {
	data      map[string]*Config
	groupData map[string]*Group // izigrava bazu podataka
}

// swagger:route POST /config/ config createConfig
// Add new config
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseConfig
func (cs *configServer) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := createId()
	rt.Id = id
	cs.data[id] = rt
	renderJSON(w, rt)
}

// swagger:route GET /configs/ config getConfigs
// Get all configs
//
// responses:
//
//	200: []ResponseConfig
func (cs *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := []*Config{}
	for _, v := range cs.data {
		allTasks = append(allTasks, v)
	}
}

// swagger:route GET /config/{id}/ config getConfigById
// Get config by ID
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseConfig
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := cs.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

// swagger:route DELETE /config/{id}/ config deleteConfig
// Delete config
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
//	201: ResponseConfig
func (cs *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := cs.data[id]; ok {
		delete(cs.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// swagger:route POST /group/ group createGroup
// Add new group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (cs *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	group, err := decodeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	group.Id = id
	cs.groupData[id] = group
	renderJSON(w, group)
}

// swagger:route PUT /group/{g_id}/config/{c_id}/ group addConfigToGroup
// Add config to group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ResponseGroup
func (cs *configServer) addConfigToGroup(w http.ResponseWriter, req *http.Request) {
	groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]
	task, ok := cs.data[id]
	group, ook := cs.groupData[groupId]
	if !ok || !ook {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	group.Configs = append(group.Configs, *task)
	cs.groupData[groupId] = group

	return
}

// swagger:route GET /groups/ group getGroups
// Get all groups
//
// responses:
//
//	200: []ResponseGroup
func (cs *configServer) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allGroups := []*Group{}
	for _, v := range cs.groupData {
		allGroups = append(allGroups, v)
	}
	renderJSON(w, allGroups)
}

// swagger:route GET /group/{id}/ group getGroupById
// Get group by ID
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseGroup
func (cs *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := cs.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

// swagger:route DELETE /group/{id}/ group deleteGroup
// Delete group
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
//	201: ResponseGroup
func (cs *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	_, ok := cs.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	delete(cs.groupData, id)
}

// swagger:route DELETE /group/{g_id}/config/{c_id}/ group deleteConfigFromGroup
// Delete config from group
//
// responses:
//
//	404: ErrorResponse
//	204: NoContentResponse
func (cs *configServer) delConfigFromGroupHandler(w http.ResponseWriter, req *http.Request) {
	groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]
	group, ok := cs.groupData[groupId]
	if !ok {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for i, config := range group.Configs {
		if config.Id == id {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			cs.groupData[groupId] = group
			return
		}
	}
	err := errors.New("config not found in group")
	http.Error(w, err.Error(), http.StatusNotFound)
	return
}

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
