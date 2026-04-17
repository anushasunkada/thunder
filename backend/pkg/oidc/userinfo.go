package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *OIDCServer) Userinfo(w http.ResponseWriter, req *http.Request) {

	//Execute if there are any Before hooks

	//build userinfo json
	//TODO instead of userprovider, we should use AuthnProvider getAttributes method to fetch user attributes based on the access token scopes
	user, err := s.userProvider.FetchUser(req.Context(), "1001")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jsonBytes, _ := json.Marshal(user)
	fmt.Println(string(jsonBytes))

	userinfo := string(jsonBytes)

	fmt.Println("Handler : Userinfo")

	//Execute if there are any After hooks that want to modify the userinfo response before it's sent back to the client
	//Example is to sign the userinfo response as a JWT or encrypt it, etc.
	for _, h := range s.userinfoHooks {
		if err := h.AfterUserinfo(req.Context(), userinfo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//TODO : change the response header type to application/jwt if the userinfo is signed as a JWT, etc.
	//TODO : return the userinfo response (signed/encrypted or plain JSON) back to the client
	w.Write(jsonBytes)
}
