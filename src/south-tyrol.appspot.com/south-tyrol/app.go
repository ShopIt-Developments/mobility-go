package purchase

import (
    "net/http"
    "github.com/julienschmidt/httprouter"

    "endpoint"
)

func init() {
    router := httprouter.New()

    o := endpoint.Car{Router: router}
    o.Router.GET("/orders", o.GetAll)
    o.Router.GET("/order/:order_id", o.GetAll)
    o.Router.POST("/order", o.Add)
    o.Router.POST("/order/accept/:order_id", o.Accept)
    o.Router.PUT("/order/:order_id", o.Update)
    o.Router.DELETE("/order/:order_id", o.Delete)

    http.Handle("/", router)
}