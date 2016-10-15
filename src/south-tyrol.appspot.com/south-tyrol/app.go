package purchase

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"endpoint"
)

func init() {
	router := httprouter.New()

	c := endpoint.Car{Router: router}
	c.Router.GET("/mobility/vehicles/available", c.GetAvailable)
	c.Router.GET("/mobility/vehicles/my", c.GetMy)
	c.Router.GET("/mobility/vehicles", c.GetAll)
	c.Router.GET("/mobility/vehicle/:vehicle_id", c.GetOne)
	c.Router.POST("/mobility/vehicle", c.Add)
	c.Router.PUT("/mobility/vehicles/:vehicle_id", c.Update)
	c.Router.DELETE("/mobility/vehicles/:vehicle_id", c.Delete)

	u := endpoint.User{Router: router}
	u.Router.GET("/mobility/user", u.Get)
	u.Router.POST("/mobility/user", u.Add)
	u.Router.PUT("/mobility/user", u.Update)
	u.Router.DELETE("/mobility/user", u.Delete)

	u.Router.GET("/mobility/points", u.GetPoints)
	u.Router.POST("/mobility/points/:points", u.AddPoints)

	r := endpoint.Rating{Router: router}
	r.Router.GET("/mobility/ratings", r.GetRatings)
	r.Router.GET("/mobility/ratings/:rating_id", r.GetOne)
	r.Router.GET("/mobility/rated", r.GetRated)
	r.Router.POST("/mobility/rating", r.Add)
	r.Router.PUT("/mobility/ratings/:rating_id", r.Update)
	r.Router.DELETE("/mobility/rating/s:rating_id", r.Delete)

	http.Handle("/", router)
}