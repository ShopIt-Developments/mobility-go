package purchase

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"endpoint"
)

func init() {
	router := httprouter.New()

	c := endpoint.Car{Router: router}
	c.Router.GET("/cars", c.GetAll)
	c.Router.GET("/car/:car_id", c.GetOne)
	c.Router.POST("/car", c.Add)
	c.Router.PUT("/car/:car_id", c.Update)
	c.Router.DELETE("/car/:car_id", c.Delete)

	u := endpoint.User{Router: router}
	u.Router.GET("/user", u.Get)
	u.Router.POST("/user", u.Add)
	u.Router.PUT("/user", u.Update)
	u.Router.DELETE("/user", u.Delete)

	http.Handle("/", router)
}