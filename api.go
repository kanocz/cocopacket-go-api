package api

import (
	"errors"
	"net"
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

// AddSlave adds slave to master on ip:port with name
// and possibly copy list of ips from just existing slave copyFrom
func AddSlave(ip net.IP, port uint16, name string, copyFrom string) error {
	var r result

	err := Send("POST", mainAPIURL+"/v1/slaves", map[string]interface{}{
		"ip":   ip.String(),
		"port": port,
		"name": name,
		"copy": copyFrom,
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

// DeleteSlave removes slave from master
func DeleteSlave(slave string) error {
	var r result

	err := Send("DELETE", mainAPIURL+"/v1/slaves?slave="+url.QueryEscape(slave), nil, &r)
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

// GroupStats returns stats for all IPs/URLs in group for about last 24 hours with 1-hour aggregation (report -> limit only to ip+slaves selected for report using frontend)
func GroupStats(group string, report bool) (GroupStatsData, error) {
	var data GroupStatsData
	reportAdd := ""
	if report {
		reportAdd = "?report=true"
	}
	err := Get(mainAPIURL+"/v1/catstats/"+url.QueryEscape(group+"->")+reportAdd, &data)
	return data, err
}
