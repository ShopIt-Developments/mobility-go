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

    http.Handle("/", router)
}