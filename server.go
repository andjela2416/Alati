package main

import (
	"encoding/json"
	"errors"
	s "example.com/mod/store"
	"fmt"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type configServer struct {
	store *s.Store
	//data      map[string]*s.Config
	groupData map[string]*s.Group // izigrava bazu podataka*/
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

	post, err := cs.store.Config(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
}

// swagger:route GET /configs/ config getConfigs
// Get all configs
//
// responses:
//
//	200: []ResponseConfig
func (cs *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := cs.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
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
	version := mux.Vars(req)["version"]
	task, err := cs.store.Get(id, version)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	version := mux.Vars(req)["version"]

	msg, err := cs.store.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

/*if v, ok := cs.data[id]; ok {
	delete(cs.data, id)
	renderJSON(w, v)
} else {
	err := errors.New("key not found")
	http.Error(w, err.Error(), http.StatusNotFound)
}*/

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

	rt, err := decodeGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := cs.store.PostGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
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
	/*groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]
	task, ok := cs.groupData[id]
	group, ook := cs.groupData[groupId]
	if !ok || !ook {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	group.Configs = append(group.Configs, *task)
	cs.groupData[groupId] = group
	*/
	groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]

	// Dekodiranje JSON podataka iz zahteva u objekat tipa Config
	var config s.Config
	err := json.NewDecoder(req.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Provera da li grupa i konfiguracija postoje
	group, ook := cs.groupData[groupId]
	if !ook {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, ok := cs.groupData[id]
	if !ok {
		err := errors.New("config not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Printf("Task ID: %s\n", task.Id)

	// Dodavanje konfiguracije u grupu
	group.Configs = append(group.Configs, config)
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
	allGroups := []*s.Group{}
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
	version := mux.Vars(req)["version"]
	task, err := cs.store.Get(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	version := mux.Vars(req)["version"]

	msg, err := cs.store.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
	/*_, ok := cs.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	delete(cs.groupData, id)*/
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

func (s *configServer) getPostByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]

	task, err := s.store.GetGroupsByLabels(id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, task)
}
