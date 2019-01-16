package api

import (
	"errors"
	"net/url"
)

var (
	mainAPIURL string
)

// Init sets API url and authorization parameters
func Init(url string, username string, password string) {
	mainAPIURL = url
	SetBasicAuth(username, password)
}

// GetConfigInfo returns current configuration
func GetConfigInfo() (ConfigInfo, error) {
	var result ConfigInfo
	err := Get(mainAPIURL+"/v1/config", &result)
	return result, err
}

// GetSlaveList returns list of defined slave probes
func GetSlaveList() ([]string, error) {
	var result map[string]string
	err := Get(mainAPIURL+"/v1/slaves", &result)
	list := make([]string, 0, len(result))
	for slave := range result {
		list = append(list, slave)
	}
	return list, err
}

// AddIP is simple interface for single IP adding
func AddIP(ip string, slaves []string, description string, groups []string, favourite bool) error {
	var r result

	err := Send("PUT", mainAPIURL+"/v1/config/ping/"+ip, TestDesc{
		Description: ip + " " + description,
		Favourite:   favourite,
		Groups:      groups,
		Slaves:      slaves,
	}, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil
}

// AddIPs function adds multiply ips using only one API call
func AddIPs(ips []string, slaves []string, description string, groups []string, favourite bool) error {
	payload := make(map[string]TestDesc, len(ips))

	for _, ip := range ips {
		payload[ip] = TestDesc{
			Description: ip + " " + description,
			Favourite:   favourite,
			Groups:      groups,
			Slaves:      slaves,
		}
	}

	var r result

	err := Send("PUT", mainAPIURL+"/v1/mconfig/add", map[string]interface{}{
		"ips": payload,
	}, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil

}

// DeleteIP removes one IP from cocopacket instance
func DeleteIP(ip string) error {
	var r result

	err := Send("DELETE", mainAPIURL+"/v1/config/ping/"+ip, nil, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil
}

// DeleteIPs function deletes multiply ips using only one API call
func DeleteIPs(ips []string) error {
	var r result

	err := Send("PUT", mainAPIURL+"/v1/mconfig/delete", map[string]interface{}{
		"ips": ips,
	}, &r)
	if nil != err {
		return err
	}

	if r.Result != "OK" && "" != r.Error {
		return errors.New(r.Error)
	}

	if r.Result != "OK" {
		return errors.New("unknown error")
	}

	return nil
}

// ListUsers return map with logins and associated boolean indicating if user is admin
func ListUsers() (map[string]bool, error) {
	var users map[string]bool
	err := Get(mainAPIURL+"/v1/users", &users)
	return users, err
}

// AddUser adds new user (or replaces existing)
func AddUser(login string, password string, admin bool) (map[string]bool, error) {

	var users map[string]bool
	t := "user"
	if admin {
		t = "admin"
	}

	err := SendForm("PUT", mainAPIURL+"/v1/users", url.Values{
		"login":  []string{login},
		"passwd": []string{password},
		"type":   []string{t},
	}, &users)
	if nil != err {
		return nil, err
	}

	return users, err
}

// DeleteUser removes user from master
func DeleteUser(login string) (map[string]bool, error) {

	var users map[string]bool

	err := Send("DELETE", mainAPIURL+"/v1/users?login="+url.QueryEscape(login), nil, &users)
	if nil != err {
		return nil, err
	}

	return users, err
}
