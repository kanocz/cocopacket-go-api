# cocopacket-go-api
cocopacket.com API implementation implemented in golang

## api examples
please look at examples folder - there are some usefull tools that are just prepared for usege covering basic functions like managing ips, users and so on

all examples has such flags:
* `-url` - link to your master instance like `http://yourname.client.cocopacket.com/`
* `-user` - login of user, usualy `admin`
* `-password` - just password

some examples has additional flags, just run to see

few examples:
```bash
./listips -url 'http://yourname.client.cocopacket.com/' -user admin -password helloWorld
./listslaves -url 'http://yourname.client.cocopacket.com/' -user admin -password helloWorld
./addslave -url 'http://yourname.client.cocopacket.com/' -user admin -password helloWorld 1.1.1.1 3030 NEWSLAVE
```

P.S.: in every examples folder just run `go build` to compile

## running cocopacket-slave from docker
some users asks if it's possible to run cocopacket-slave in a docker containter. Sure it's possible, but you'll need some basic knowlage about docker itself and port forwarding. Example docker file:

```docker
FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD cocopacket-slave /
ADD slave.conf /
CMD ["/cocopacket-slave", "-config", "slave.conf"]
```

create `slave.conf` like this (please adjust `aesKey` to be similar as you have on master):
```
{
  "listen": "0.0.0.0:3030",
  "memory": 3600,
  "source": "0.0.0.0",
  "aesKey": "00112233445566778899AABBCCDDEEFF"
}
```

please make sure that you have installed `jq` utility used to parse json to obtain last version. download and build:

```bash
cp /etc/ssl/certs/ca-certificates.crt . &&
curl cocopacket-slave "https://updates.cocopacket.cloud/cocopacket-slave/`curl -s https://updates.cocopacket.cloud/cocopacket-slave/linux-amd64.json  | jq -r .Version`/linux-amd64.gz" | gzip -dc > cocopacket-slave &&
chmod 755 cocopacket-slave &&
docker build -t cocoslave .
```

and then run with something like this (adjust ports):

```bash
docker run --publish 3030:3030 --name coco1 cocoslave

```