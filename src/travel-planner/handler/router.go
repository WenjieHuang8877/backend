package handler

import (
	"net/http"
	"travel-planner/util"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var mySigningKey []byte

func InitRouter(config *util.TokenInfo) http.Handler {
	mySigningKey = []byte(config.Secret)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	router := mux.NewRouter()

	router.Handle("/app/{id}", jwtMiddleware.Handler(http.HandlerFunc(ExampleHandler))).Methods("DELETE")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(ExampleHandler)).Methods("POST")

	router.Handle("/user/signin", http.HandlerFunc(SigninHandler)).Methods("POST")
	router.Handle("/user/{id}",jwtMiddleware.Handler(http.HandlerFunc(UpdateUserHander))).Methods("POST")
	router.Handle("/user/getUser/{id}",jwtMiddleware.Handler(http.HandlerFunc(GetUserHandler))).Methods("GET")
	router.Handle("/vacation/MyVacation",jwtMiddleware.Handler(http.HandlerFunc(GetSitesHandler))).Methods("GET")
	router.Handle("/vacation",jwtMiddleware.Handler(http.HandlerFunc( SearchSitesHandler))).Methods("POST")
  router.Handle("/vacation/{vacation_id}/sites/{id}", jwtMiddleware.Handler(http.HandlerFunc(addSiteHandler))).Methods("POST")


	// TODO: add jwtMiddleware.Handler() wrapper
	router.Handle("/vacation", jwtMiddleware.Handler(http.HandlerFunc(GetVacationsHandler))).Methods("GET")
	router.Handle("/vacation/init", jwtMiddleware.Handler(http.HandlerFunc(SaveVacationsHandler))).Methods("POST")

	router.Handle("/vacation/{vacation_id}/plan", http.HandlerFunc(GetVacationPlanHandler)).Methods("GET")
	router.Handle("/vacation/{vacation_id}/plan/init", http.HandlerFunc(InitPlanHandler)).Methods("POST")
	router.Handle("/vacation/{vacation_id}/plan/{plan_id}/save", http.HandlerFunc(SaveActivitiesHandler)).Methods("POST")
  
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
