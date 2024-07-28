package user

import (
	"context"

	"net/http"
	"reflect"

	v "github.com/core-go/core/v10"
	mgo "github.com/core-go/mongo"
	"github.com/core-go/search"
	"github.com/core-go/search/mongo/query"

	"go.mongodb.org/mongo-driver/mongo"

	"go-service/internal/user/handler"
	"go-service/internal/user/model"
	"go-service/internal/user/service"
)

type UserTransport interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(db *mongo.Database, logError func(context.Context, string, ...map[string]interface{})) (UserTransport, error) {
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userType := reflect.TypeOf(model.User{})
	userQuery := query.UseQuery(userType)
	userSearchBuilder := mgo.NewSearchBuilder(db, "users", userQuery, search.GetSort)
	userRepository := mgo.NewRepository(db, "users", userType)
	userService := service.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userSearchBuilder.Search, userService, logError, validator.Validate, nil)
	return userHandler, nil
}
