package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

type Sessions struct {
	*revel.Controller
}

func (s *Sessions) Login() revel.Result {
	//The user auth struct that we will use to decode the json into
	type userAuth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
	auth := userAuth{}
	//get the info from the auth request
	err := DecodeJsonRequest(s.Controller, &auth)
	if _err, _result := CheckStopCondition(s.Controller, 500, err); _err {
		return _result
	}
	if auth.User != "" && auth.Password != "" {
		user, err := models.GetUserByEmailOrUsername(auth.User)
		if _err, _result := CheckStopCondition(s.Controller, 500, err); _err {
			return _result
		}

		hashedPassword := user.GetHashedPassword()
		//this means we didn't match anything
		if hashedPassword == nil {
			return ReturnLast(CheckStopCondition(s.Controller, 401, "")).(revel.Result)
		}

		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(auth.Password))
		if err == nil {
			s.Session["user"] = user.GetId()
			s.Session["userType"] = user.GetType()
			result, _ := user.Render()
			return result
		}
	}

	return ReturnLast(CheckStopCondition(s.Controller, 401, "")).(revel.Result)

}

func (s *Sessions) Logout() revel.Result {
	s.Session = map[string]string{}
	return ReturnLast(CheckStopCondition(s.Controller, 200, "")).(revel.Result)
}
