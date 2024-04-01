# Rate limiter

A test task for a position Go Middle Developer

## Description

Imagine we have a service which is receiving a huge number of requests,
but it can only serve a limited number of requests per specific period of time.
To handle this problem we would need some kind of throttling or rate limiting mechanism that would allow
only a certain number of requests so our service can respond to all of them.
A rate limiter, at a high-level, limits the number of events an entity (user, device, IP, etc.)
can perform in a particular time window.

For example:

- A user can send only one message per second.
- A user is allowed only three failed credit card transactions per day.
- A single IP can only create twenty accounts per day.

In general, a rate limiter caps how many requests a sender can issue in a specific time window.
It then blocks requests once the cap is reached.

### Requirements

It should be built using Golang or Node.Js or Rust.
Repository could contain Readme with instructions how to run and test application.
Test coverage will be plus.
Upload your solution on Github and send us a link to repository.
Don't spend more than 4 hours, it's a simple task and shouldn't be over engineered.

## Solution

### How to run

#### Prerequisites

I expect you to have installed Docker and (optionally) CMake

#### With CMake

##### Start

Just run following command

```shell
make
```

It will start server and any required infrastructure and will show you server logs.

##### Restart

Run this again

```shell
make
```

It will rebuild server binary and restart server (in case of changes).

###### Force restart

If you need to force container recreation

```shell
make recreate
```

##### Stop
This command will stop and remove all the containers
```shell
make stop
```

#### Without CMake

```shell
docker compose -f 
```