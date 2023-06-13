package main

import (
	"context"
	"errors"
	s "example.com/mod/store"
	tracer "example.com/mod/tracer"
	"fmt"
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

func NewPostServer() (*configServer, error) {
	store, err := s.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(tracer)
	return &configServer{
		store:  store,
		tracer: tracer,
		closer: closer,
	}, nil
}
func (s *configServer) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *configServer) GetCloser() io.Closer {
	return s.closer
}

func (s *configServer) CloseTracer() error {
	return s.closer.Close()
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
	span := tracer.StartSpanFromRequest("creteConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling config create at %s\n", req.URL.Path)),
	)

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
	ctx := tracer.ContextWithSpan(context.Background(), span)
	rt, err := decodeBody(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//post, err := cs.store.Config(rt)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if cs.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}
	post, err := cs.store.Config(ctx, rt)

	reqId := ""

	if err == nil {
		reqId = cs.store.SaveRequestId(ctx)
	}

	renderJSON(ctx, w, post)
	renderJSON(ctx, w, "Idempotence key:"+reqId)
}

// swagger:route GET /configs/ config getConfigs
// Get all configs
//
// responses:
//
//	200: []ResponseConfig
func (cs *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getAllConfigsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all configs at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	allTasks, err := cs.store.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, allTasks)
}

// swagger:route GET /config/{id}/ config getConfigById
// Get config by ID
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseConfig
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all configs at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	task, err := cs.store.Get(ctx, id, version)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, task)
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
	span := tracer.StartSpanFromRequest("deleteConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete config at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]

	version := mux.Vars(req)["version"]

	msg, err := cs.store.Delete(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, msg)
}
func (cs *configServer) delConfigByLabelHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("deleteConfigHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling delete config at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]

	version := mux.Vars(req)["version"]

	label := mux.Vars(req)["labels"]

	msg, err := cs.store.DeleteByLabel(ctx, id, version, label)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, msg)
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
	span := tracer.StartSpanFromRequest("creteGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling group create at %s\n", req.URL.Path)),
	)
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
	ctx := tracer.ContextWithSpan(context.Background(), span)
	rt, err := decodeGroup(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//post, err := cs.store.PostGroup(rt)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if cs.store.FindRequestId(ctx, requestId) == true {
		http.Error(w, "Request has been already sent", http.StatusBadRequest)
		return
	}
	post, err := cs.store.PostGroup(ctx, rt)

	reqId := ""

	if err == nil {
		reqId = cs.store.SaveRequestId(ctx)
	}

	renderJSON(ctx, w, post)
	renderJSON(ctx, w, "Idempotence key:"+reqId)
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
	span := tracer.StartSpanFromRequest("addConfigToGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling add config to group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	groupId := mux.Vars(req)["g_id"]
	groupVersion := mux.Vars(req)["g_version"]
	id := mux.Vars(req)["c_id"]
	configVersion := mux.Vars(req)["c_version"]

	group, err := cs.store.GetOneGroup(ctx, groupId, groupVersion)
	if err != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, err := cs.store.GetOneConfig(ctx, id, configVersion)
	if err != nil {
		err := errors.New("config not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	group.Configs = append(group.Configs, *task)

	cs.store.SaveGroup(ctx, group)
	renderJSON(ctx, w, group)
}
func (cs *configServer) addConfigToGroup2(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigToGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling add config to group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	groupId := mux.Vars(req)["g_id"]
	id := mux.Vars(req)["c_id"]

	group, err := cs.store.GetOneGroup2(ctx, groupId)
	if err != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, err := cs.store.GetOneConfig2(ctx, id)
	if err != nil {
		err := errors.New("config not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	group.Configs = append(group.Configs, *task)

	cs.store.SaveGroup(ctx, group)
	renderJSON(ctx, w, group)
}

// swagger:route GET /groups/ group getGroups
// Get all groups
//
// responses:
//
//	200: []ResponseGroup
func (cs *configServer) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getAllGroupsHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all groups at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	allTasks, err := cs.store.GetAllGroups(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, allTasks)
}

// swagger:route GET /group/{id}/ group getGroupById
// Get group by ID
//
// responses:
//
//	404: ErrorResponse
//	200: ResponseGroup
func (cs *configServer) getGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all groups at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	task, err := cs.store.GetGroup(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, task)

}

func (cs *configServer) getGroupHandlerId(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all groups at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]

	task, err := cs.store.GetGroupId(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, task)

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
	span := tracer.StartSpanFromRequest("delConfigFromGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config from group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	msg, err := cs.store.DeleteGroup(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, msg)
	/*_, ok := cs.groupData[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	delete(cs.groupData, id)*/
}
func (cs *configServer) delGroupHandlerId(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delConfigFromGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config from group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]

	msg, err := cs.store.DeleteGroupId(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, msg)
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

	span := tracer.StartSpanFromRequest("delConfigFromGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config from group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	groupId := mux.Vars(req)["groupId"]
	groupVersion := mux.Vars(req)["g_version"]
	//configVersion := mux.Vars(req)["c_version"]
	id := mux.Vars(req)["id"]
	group, err2 := cs.store.GetOneGroup(ctx, groupId, groupVersion)
	if err2 != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for i, config := range group.Configs {
		if config.Id == id {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			//cs.groupData[groupId] = group
			grupas, err := cs.store.SaveGroup(ctx, group)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			renderJSON(ctx, w, grupas)
			return
		}
	}
	err := errors.New("config not found in group")
	http.Error(w, err.Error(), http.StatusNotFound)
	grupas, err := cs.store.SaveGroup(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, grupas)
}
func (cs *configServer) delConfigFromGroupHandler2(w http.ResponseWriter, req *http.Request) {

	span := tracer.StartSpanFromRequest("delConfigFromGroupHandler", cs.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del config from group at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	groupId := mux.Vars(req)["groupId"]
	//groupVersion := mux.Vars(req)["g_version"]
	//configVersion := mux.Vars(req)["c_version"]
	id := mux.Vars(req)["id"]
	group, err2 := cs.store.GetOneGroup2(ctx, groupId)
	if err2 != nil {
		err := errors.New("group not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for i, config := range group.Configs {
		if config.Id == id {
			group.Configs = append(group.Configs[:i], group.Configs[i+1:]...)
			//cs.groupData[groupId] = group
			grupas, err := cs.store.SaveGroup(ctx, group)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			renderJSON(ctx, w, grupas)
			return
		}
	}
	err := errors.New("config not found in group")
	http.Error(w, err.Error(), http.StatusNotFound)
	grupas, err := cs.store.SaveGroup(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, grupas)
}

func (ts *configServer) swaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}

func (s *configServer) getPostByLabel(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getPostByLabelHandler", s.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get post by label at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]

	task, err := s.store.GetConfigsByLabels(ctx, id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, task)
}

/*
func (s *configServer) getGroupsByLabel(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupsByLabelHandler", s.tracer, req)
	defer span.Finish()

	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get groups by label at %s\n", req.URL.Path)),
	)

	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]

	task, err := s.store.GetGroupsByLabels(ctx, id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, task)
}*/
