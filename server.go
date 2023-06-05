package main

import (
	"errors"
	s "example.com/mod/store"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"io"
	"mime"
	"net/http"
)

const (
	name = "post_service"
)

type configServer struct {
	store  *s.Store
	tracer opentracing.Tracer
	closer io.Closer
	//data      map[string]*s.Config
	//groupData map[string]*s.Group
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
	requestId := req.Header.Get("x-idempotency-key")
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

	//post, err := cs.store.Config(rt)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if cs.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}
	post, err := cs.store.Config(rt)

	reqId := ""

	if err == nil {
		reqId = cs.store.SaveRequestId()
	}

	renderJSON(w, post)
	renderJSON(w, "Idempotence key:"+reqId)
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
	requestId := req.Header.Get("x-idempotency-key")
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

	//post, err := cs.store.PostGroup(rt)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if cs.store.FindRequestId(requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}
	post, err := cs.store.PostGroup(rt)

	reqId := ""

	if err == nil {
		reqId = cs.store.SaveRequestId()
	}

	renderJSON(w, post)
	renderJSON(w, "Idempotence key:"+reqId)
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
	groupVersion := mux.Vars(req)["g_version"]
	id := mux.Vars(req)["c_id"]
	configVersion := mux.Vars(req)["c_version"]
	/*group, groupExists := cs.groupData[groupId]
	if !groupExists {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}*/

	group2, err := cs.store.GetOneGroup(groupId, groupVersion)
	if err != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, err := cs.store.GetOneConfig(id, configVersion)
	if err != nil {
		err := errors.New("config not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	/*config := s.Config{
		Id: task.Id,
	}*/
	group2.Configs = append(group2.Configs, *task)
	//cs.groupData[groupId] = group2
	/*grupas, err := cs.store.SaveGroup(group2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/
	cs.store.SaveGroup(group2)
	renderJSON(w, group2)
	/*groupId := mux.Vars(req)["g_id"]
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

	return*/
}

// swagger:route GET /groups/ group getGroups
// Get all groups
//
// responses:
//
//	200: []ResponseGroup
func (cs *configServer) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {

	allTasks, err := cs.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
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
	task, err := cs.store.GetGroup(id, version)
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

	msg, err := cs.store.DeleteGroup(id, version)
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
	/*groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]
	group, ok := cs.groupData[groupId]
	if !ok {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}*/
	groupId := mux.Vars(req)["groupId"]
	groupVersion := mux.Vars(req)["g_version"]
	//configVersion := mux.Vars(req)["c_version"]
	id := mux.Vars(req)["id"]
	group, err2 := cs.store.GetOneGroup(groupId, groupVersion)
	if err2 != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for i, config := range group.Configs {
		if config.Id == id {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			//cs.groupData[groupId] = group
			grupas, err := cs.store.SaveGroup(group)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			renderJSON(w, grupas)
			return
		}
	}
	err := errors.New("config not found in group")
	http.Error(w, err.Error(), http.StatusNotFound)
	grupas, err := cs.store.SaveGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, grupas)
}

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}

func (s *configServer) getPostByLabel(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]

	task, err := s.store.GetConfigsByLabels(id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, task)
}

func (s *configServer) getGroupsByLabel(w http.ResponseWriter, req *http.Request) {
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
