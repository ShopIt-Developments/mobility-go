package purchase

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"endpoint"
)

func init() {
	router := httprouter.New()

    v := endpoint.Vehicle{Router: router}
    v.Router.GET("/mobility/vehicles/available", v.GetAvailable)
    v.Router.GET("/mobility/vehicles/booked", v.GetBooked)
    v.Router.GET("/mobility/vehicles/my", v.GetMy)
    v.Router.GET("/mobility/vehicles", v.GetAll)
    v.Router.GET("/mobility/vehicle/:vehicle_id", v.GetOne)
    v.Router.POST("/mobility/vehicle", v.New)
    v.Router.PUT("/mobility/vehicles/:vehicle_id", v.Update)
    v.Router.DELETE("/mobility/vehicle/:vehicle_id", v.Delete)

    o := endpoint.Order{Router: router}
    o.Router.POST("/mobility/order/:vehicle_id", o.New)
    o.Router.DELETE("/mobility/order/:order_id", o.Delete)

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
	r.Router.DELETE("/mobility/rating/:rating_id", r.Delete)

    t := endpoint.Trip{Router: router}
    t.Router.POST("/mobility/trip", t.New)

	p := endpoint.Payments{Router: router}
	p.Router.POST("/mobility/payments/scan/:order_id", p.Scan)
	p.Router.POST("/mobility/payments/accept/:order_id", p.Accept)

	http.Handle("/", router)
}