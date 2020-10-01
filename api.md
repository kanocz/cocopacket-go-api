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
* `slavesThresholds` object(map) with string keys (slave names) and 3-field object as a value containig `lossThreshold`, `latencyThreshold` and `timeThreshold` overrides for specified slave(s) ![1.0.4-7](https://img.shields.io/static/v1?label=ver&message=1.0.4-7&color=white)

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

## set maintenance list ![1.0.3-6](https://img.shields.io/static/v1?label=ver&message=1.0.3-6&color=white)
*PUT* `/v1/maintenance` _admin_  
json-payload:
```json
[ "8.8.8.8", "1.1.1.1" ]
```
set list of IPs for which push notifications is *off*

## get current maintenance list ![1.0.3-6](https://img.shields.io/static/v1?label=ver&message=1.0.3-6&color=white)
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

## configure push notifications ![1.0.2-0](https://img.shields.io/static/v1?label=ver&message=1.0.2-0&color=white)
*PUT* `/v1/notify` _admin_  
put whole list of push notify destinations
```json
{
    "sms": {
	    "method": "POST",
	    "url": "https://sms.gateway.com/something-more..",
	    "payload": "{ some json... }",
	    "contentType": "application/json",
	    "headers": {
            "Auth": "token"
        },
	    "frequency": 60,
	    "frequencyPerIP": 300,
	    "minSlavesFailed": 1
    },
    "telegram": ...
}
```
`method` GET/PUT/POST  
`url` just a url :) 
`payload`  what to send, valid for PUT/POST  
`contentType`  usualy application/json or text/plain  
`headers`  additional HTTP headers needed
`frequency`  one notification per "frequency" secodns (for example 60 means that only one message per minute will be send), can be set to zero  
`frequencyPerIP`  one notification per "frequency" secodns for one ip (for example 300 means that only one message per 5 minutes will be send for failed ip... but other messages will be send regarding to `frequency`)  
`minSlavesFailed` ![1.0.4-6](https://img.shields.io/static/v1?label=ver&message=1.0.4-6&color=white) send message only if at least minSlavesFailed slaves reports a problem with one IP   

in `url` and `payload` any `<*IP*>`, `<*SLAVE*>`, `<*LATENCY*>`, `<*LOSS*>` and `<*GROUP*>` will be replaced with current inident values  
in case of `minSlavesFailed > 1` data from message received from last slave that triggered incident is used  
additional values `<*SCOUNT*>` (at least `<*SCOUNT*>` slaves triggered) and `<*SLAVES*>` (comma separated list of slave at the moment of event pushing) added ![1.0.4-6](https://img.shields.io/static/v1?label=ver&message=1.0.4-6&color=white)

push notifications are processes over 60-second slices of data so it's always be about 1-minute delay

## get current list of push notification configurations ![1.0.2-0](https://img.shields.io/static/v1?label=ver&message=1.0.2-0&color=white)
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

## get current password recovery mail template ![1.0.4-4](https://img.shields.io/static/v1?label=ver&message=1.0.4-4&color=white)

*GET* `/v1/config/preset`
i

## set new password recovery mail template ![1.0.4-4](https://img.shields.io/static/v1?label=ver&message=1.0.4-4&color=white)

*POST* `/v1/config/preset`
json-payload:
```json
{
    "subject": "CocoPacket password reset",
    "contentType": "text/plain",
    "body": "Your password for CocoPacket on address {URL} is reseted to {PASSWORD}"
}
```
`{URL}` and `{PASSWORD}` will be replaced with actual URL of instance and new password  
in case of empty subject/contentType/body in *PUT* call they will be filled with default values
