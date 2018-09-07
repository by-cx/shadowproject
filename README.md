# ShadowProject

This project makes easier to run a single application in cloud.

## Problem

It's expensive to run Python, Go or Node.js applications, because you need
your own server and that means a lot of time maitaining it or you have
to pay somebody. This project wants to make this easy and cheap with a simple
principle. Your application doesn't need the server all the time but only
in moments the requests is coming. Any other time it can be turned off.
That means you won't pay $10 per app but just $1 per app. The price depends
on the storage and time the application is running, but for many use cases
it will be always cheaper.

The technology behind this is in DigitalOcean's and AWS's block storage
services (EBS in AWS). This tools holds your data in this reliable storage
and mounts it only in case the requests comes. When the app is not used
everything related except the data is removed and possibly even the
underlaying server too. That means the most of the time we pay for
storage but not for servers.  
 
 
## Architecture

Every application is handled internally as a task. If request comes
the task is loaded from the databases and the environment is preapared.

...

## Longterm TODO:

* Support for stateful tasks
* Support for stateless tasks - no permanent storage
* Autoscaling for stateless tasks based on CPU load
