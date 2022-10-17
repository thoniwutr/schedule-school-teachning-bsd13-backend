package handlers

import (
	"fmt"
	"github.com/schedule-school-teachning-bsd13-backend/util"
	"net/http"
)


type ScheduleHandler struct {
	util       *util.HandlerUtil
}

func NewScheduleHandler(util *util.HandlerUtil) *ScheduleHandler {
	return &ScheduleHandler{util: util}
}

func (m *ScheduleHandler) ScheduleClass(rw http.ResponseWriter, r *http.Request) {


	membershipNumber := r.FormValue("memberkey")
	fmt.Println("GET params were:", membershipNumber)

	m.util.HTTPSuccess(rw, "success", http.StatusCreated)
}


