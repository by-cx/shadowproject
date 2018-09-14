# ShadowProject

This project makes easier to run a single application in cloud.

It's under development, but getting better every day. Right even the
API and data structure are not complete so don't use this in production.

## Problem

It's expensive to run Python, Go or Node.js applications, because you need
your own server and that means a lot of time to maitaining it or money to pay
somebody to do that for you. This project wants to make this easy and cheap.
Your application doesn't need the server all the time but only in moments
the requests is coming. Any other moment the app can be turned off. That
means you won't pay $10 per app each month but for example $1 per app for
same time. The price depends on the storage and time the
application is running. For many use cases it will be always cheaper than
the standard hosting or VPS.

The technology behind this is in DigitalOcean's and AWS's block storage
services (EBS in AWS). This tools holds your data in this reliable storage
and this tools mounts it only in case the request comes. When the app is not used
everything related, except the data, is removed and possibly the
underlaying server too. That means the most of the time you pay for
storage but not for servers.  

## Architecture

Every application is handled internally as a task. The request goes to main proxy
where it's redirected to a node proxy. Node proxy cares only about local
containers and it fires the container for the coming request regardless
anything else. The proxy server remembers where the request was
redirected and another request will be sent to the same node. Node decides
what containers to kill and what let to run.

The architecture is not finished yet and a lot of things can change.

![Shadow scheme](contrib/Shadow.jpg)

## Longterm TODO:

* Support for stateful tasks
* Support for stateless tasks - no permanent storage
* Autoscaling for stateless tasks based on CPU load
* Accounting