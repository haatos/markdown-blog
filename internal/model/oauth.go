package model

import "strings"

type GoogleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarURL string `json:"picture"`
}

type GithubUser struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func (gu *GithubUser) GetFirstNameLastName() (string, string) {
	split := strings.Split(gu.Name, " ")
	if len(split) == 1 {
		return gu.Name, ""
	}
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], split[len(split)-1]
}
