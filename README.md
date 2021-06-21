# socks5-callback-server

The purpose of this server is to enable quickly create Internet servers
which allow software to automatically check if a public port
needed for operation is open from the Internet.

For example, what if Caddy had an option `-check-public-ports` that told users if the 
necessary ports are open or not when running an HTTP server.
This can dramatically simplify the life of end-users who are not good at using Firewalls, port-forwarding, `telnet` and `curl` for checking whether TCP services are properly
presenting outside the firewall.

This server is restricted to limit it's usefulness: it only permits proxying TCP to same IP address as request is coming from. It really can't be used for much beyond checking if my firewall has
a port forwarded, and there is a listening TCP service on that port.

You would general run this server on an affordable $5 or $10 virtual machine, that has a direct routable IP address on it's interfaces.

It's Dockerfied to make it simple to get going.

Default port is 60000

## Docker build and rust run


Build:
```bash
docker build -t x186k/socks5-callback-server .
```

For non-detached testing:
```bash
docker run --network host x186k/socks5-callback-server
```

## Upload to dockerhub:

```bash
docker login
docker push x186k/socks5-callback-server:latest
```

## Run using Docker in production 

```
ufw allow 60000/tcp
```

For 24x7 service on a cloud provider instance:
```bash
docker run --network host -d --restart unless-stopped x186k/socks5-callback-server
```

## Testing

Test proxy is hittable.
Run from outside proxy.
```
# replace 1.1.1.1 with proxy ip
# this should connect to proxy
curl -v telnet://1.1.1.1:60000
```

Make sure you can't use proxy as general proxy.
Run from anywhere.
```
# replace 1.1.1.1 with proxy ip
# this should fail to connect to google
curl -x socks5://1.1.1.1:60000 -v http://google.com:80
```

Test the proxy for it's intended function.
Run from outside proxy, on a box with 22/ssh open.
```
# replace 1.1.1.1 with proxy ip
# replace 2.2.2.2 with ip of box you are testing from, with 22 open
# this should connect
curl -x socks5://1.1.1.1:60000 -v telnet://1.1.1.1:22
```


