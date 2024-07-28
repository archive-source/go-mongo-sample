package handler

import (
	"context"
	"net/http"
	"reflect"

	"github.com/core-go/core"
	"github.com/core-go/search"

	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, int64) (int64, error), service service.UserService, logError func(context.Context, string, ...map[string]interface{}), validate func(context.Context, interface{}) ([]core.ErrorMessage, error), action *core.ActionConfig) *UserHandler {
	filterType := reflect.TypeOf(model.UserFilter{})
	modelType := reflect.TypeOf(model.User{})
	params := core.CreateParams(modelType, logError, validate, action)
	searchHandler := search.NewSearchHandler(find, modelType, filterType, logError, params.Log)
	return &UserHandler{service: service, SearchHandler: searchHandler, Params: params}
}

type UserHandler struct {
	service service.UserService
	*search.SearchHandler
	*core.Params
}

func (h *UserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Load(r.Context(), id)
		core.Return(w, r, res, err, h.Error, nil)
	}
}
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := core.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &user)
			core.AfterCreated(w, r, &user, res, er3, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user model.User
	er1 := core.DecodeAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Update) {
			res, er3 := h.service.Update(r.Context(), &user)
			core.HandleResult(w, r, &user, res, er3, h.Error, h.Log, h.Resource, h.Action.Update)
		}
	}
}
func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user model.User
	r, json, er1 := core.BuildMapAndCheckId(w, r, &user, h.Keys, h.Indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !core.HasError(w, r, errors, er2, h.Error, h.Log, h.Resource, h.Action.Patch) {
			res, er3 := h.service.Patch(r.Context(), json)
			core.HandleResult(w, r, json, res, er3, h.Error, h.Log, h.Resource, h.Action.Patch)
		}
	}
}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := core.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Delete(r.Context(), id)
		core.HandleDelete(w, r, res, err, h.Error, h.Log, h.Resource, h.Action.Delete)
	}
}
