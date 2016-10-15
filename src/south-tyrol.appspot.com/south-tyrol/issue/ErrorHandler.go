package issue

import (
	"net/http"
	"encoding/json"
	"model"
)

func Handle(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		e, _ := json.Marshal(model.Error{Error: err.Error()})

		w.WriteHeader(statusCode)
		w.Write(e)
	}
}
