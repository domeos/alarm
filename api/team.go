package api

import (
	"fmt"
	"github.com/domeos/alarm/g"
	"github.com/toolkits/container/set"
	"github.com/toolkits/net/httplib"
	"log"
	"strings"
	"sync"
	"time"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UsersWrap struct {
	Msg   string  `json:"msg"`
	Users []*User `json:"users"`
}

type UsersCache struct {
	sync.RWMutex
	M map[string][]*User
}

var Users = &UsersCache{M: make(map[string][]*User)}

func (this *UsersCache) Get(team string) []*User {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[team]
	if !exists {
		return nil
	}

	return val
}

func (this *UsersCache) Set(team string, users []*User) {
	this.Lock()
	defer this.Unlock()
	this.M[team] = users
}

func UsersOf(team string) []*User {
	users := CurlTeam(team)

	if users != nil {
		Users.Set(team, users)
	} else {
		users = Users.Get(team)
	}

	return users
}

func GetUsers(teams string) map[string]*User {
	userMap := make(map[string]*User)
	arr := strings.Split(teams, ",")
	for _, team := range arr {
		if team == "" {
			continue
		}

		users := UsersOf(team)
		if users == nil {
			continue
		}

		for _, user := range users {
			userMap[user.Name] = user
		}
	}
	return userMap
}

// return phones, emails
func ParseTeams(teams string) ([]string, []string) {
	if teams == "" {
		return []string{}, []string{}
	}

	userMap := GetUsers(teams)
	phoneSet := set.NewStringSet()
	mailSet := set.NewStringSet()
	for _, user := range userMap {
		if (user.Phone != "") {
			phoneSet.Add(user.Phone)
		}
		if (user.Email != "") {
			mailSet.Add(user.Email)
		}
	}
	return phoneSet.ToSlice(), mailSet.ToSlice()
}

func CurlTeam(team string) []*User {
	if team == "" {
		return []*User{}
	}

	uri := fmt.Sprintf("%s/api/alarm/group/users/wrap/", g.Config().Api.DomeOS)
	req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
	req.Param("group", team)

	var usersWrap UsersWrap
	err := req.ToJson(&usersWrap)
	if err != nil {
		log.Printf("curl %s fail: %v", uri, err)
		return nil
	}

	if usersWrap.Msg != "" {
		log.Printf("curl %s return msg: %v", uri, usersWrap.Msg)
		return nil
	}

	return usersWrap.Users
}