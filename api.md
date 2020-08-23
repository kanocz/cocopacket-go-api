# subset of API for config manipulation

## replace whole "config" (list of test targets)
*PUT* `/v1/config` _admin_  
except JSON payload, same format as retrived with GET  

## get current configuration
*GET* `/v1/config`   

## remove one entry (ip or http test)
*DELETE* `/v1/config/:type/:ip` _admin_  
`:type` can be `ping` or `http` in stable version  

## add test or replace test configuration
*PUT* `/v1/config/:type/:ip` _admin_|_user_  
access to this API can be configured to just admin or to any user  
`:type` can be `ping` or `http` in stable version  
JSON-payload represents test params:
```json 
{
	"cat":    ["GROUP->SUBGROUP->", "OTHER->"],
	"desc":   "something about IP or HTTP test",
	"fav":    false,
	"slaves": ["PRAGUE", "LONDON"],
	"report": [],
	"as":     0,
	"expire": "2020-12-01T23:50:00Z"
}
```
*warning*: in case of non-zero (`0001-01-01T00:00:00Z`) expire value ip/test will be auto-removed from system at specified date. Used for temporary add IPs for example for auto-added via support tickets.  
`as` is filled on backend side if leaved 0  

## remove group (category)
*DELETE* `/v1/group/:group` _admin_  
*warning* all IPs/tests in group which not present in another group will be deleted  
operation is not recursive, so deletion of "GROUP->" will not delete "GROUP->SUBGROUP->"  

## get configuration for group
*GET* `/v1/group/:group` _admin_  

## set configuration for group
*POST* `/v1/group/:group` _admin_  
JSON-payload have same format as GET has  
fields description:
* `isPublic`  export group to public URL
* `lossThreshold`, `latencyThreshold` and `timeThreshold` are for push notifications
* `pushNotifyA` string list of push notify configurations applied to this group
* `isAutoGroup` if `true` then backned is trying to monitor that it's enough live IPs in this group and in worst case auto add them
* `agNetwork` network (like 8.8.8.0/24) or AS (like AS15169) where too lookup live IPs for this group
* `agCount` how many live IPs this group needs
* `agSlaves` on which slaves add new IPs

## rename group
*PUT* `/v1/group/:group` _admin_  
json-payload
```json
{ "newname": "NEW->" }
```

## update slaves for all IPs/tests in group
*PUT* `/v1/groupslaves/:group` _admin_  
json-payload
```json
{
  "slaves": {
      "PRAGUE": true,
      "LONDON": false
  }
  "recursive": true
}
```
all unspecified slaves will remain unchanged  
slaves with `true` value will be added to all IPs/tests in group (and subgroups in case of recursive)  
slaves with `false` values will be removed from all IPs/tests in group (and subgroups in case of recursive)  

## get hash for public ip link
*GET* `/v1/linkkey/:ip`   

## set maintenance list
*PUT* `/v1/maintenance` _admin_  
json-payload:
```json
[ "8.8.8.8", "1.1.1.1" ]
```
set list of IPs for which push notifications is *off*

## get current maintenance list
*GET* `/v1/maintenance`   

## add multiply IPs/tests at the same time
*PUT* `/v1/mconfig/add` _admin_|_user_  
json-payload
```json
{
    "1.1.1.1": {
        "cat":    ["GROUP->SUBGROUP->", "OTHER->"],
	    "desc":   "something about IP or HTTP test",
        "fav":    false,
        "slaves": ["PRAGUE", "LONDON"],
        "report": [],
        "as":     0,
        "expire": "2020-12-01T23:50:00Z"
    },
    "8.8.8.8": {
        "cat":    ["GROUP->SUBGROUP->", "OTHER->"],
	    "desc":   "google dns",
        "fav":    false,
        "slaves": ["PRAGUE", "LONDON", "PARIS"],
        "report": [],
        "as":     0,
        "expire": "0001-01-01T00:00:00Z"
    }
}
```

## delete many IPs at the same time
*PUT* `/v1/mconfig/delete` _admin_  
json-payload
```json
{
    "ips": ["8.8.8.8", "1.1.1.1"]
}
```

## add/remove slaves for list of ips
*PUT* `/v1/mconfig/slaves` _admin_  
json-payload
```json
{
    "ips": ["8.8.8.8", "1.1.1.1"],
    "slaves": {
      "PRAGUE": true,
      "LONDON": false
    }
}
```
rules is the same as for group-slave-add/remove

## configure push notifications
*PUT* `/v1/notify` _admin_  

## get current list of push notification configurations
*GET* `/v1/notify`   

## send test notification
*GET* `/v1/notify/:type/:ip/:group?slave=XXX?&avg=10.20&loss=5&push=pushConfigurationName` _admin_  
just send testing push notification, usable for debuginh

## add / update slave
*POST* `/v1/slaves` _admin_  
json-payload:
```json
{
	"ip":    "1.1.1.1",
	"port":  1234,
	"name":  "NEW-SLAVE",
	"copy":  "LONDON",
	"update": false,
}
```
if `copy` is not empty then add to new created slave all ips/tests exist on specified existing slave  
if `update` is `true` then just update IP/port of just existing slave  

## get list of configured slaves
*GET* `/v1/slaves`   

## remove slave from system
*DELETE* `/v1/slaves?slave=XXX` _admin_  
*warning* also removes all IPs that have not other slaves

## rename slave
*PUT* `/v1/slaves` _admin_  
json-payload:
```json
{
    "oldName": "PRG-DEFAULT",
    "newName": "PRG-COGENT"
}
```

## set slaves list for ip
*PUT* `/v1/slaves/:ip` _admin_  
json-payload:
```json
   ["PRAGUE", "LONDON"]
```

## get slaves list for ip/test
*GET* `/v1/slaves/:ip`   

## get slaves status
*GET* `/v1/status/slaves`   

## get backend version and license info
*GET* `/v1/status/version`   
