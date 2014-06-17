# Hijacking Experiment

Ever want a TCP connection to a dyno?

Here's an example of doing so in go.

1. First you need to have [websockets](https://devcenter.heroku.com/articles/heroku-labs-websockets) enabled.
1. heroku create -b https://github.com/kr/heroku-buildpack-go.git \<app name\>
1. git push heroku master
1. APP=\<app name\> go run client/main.go


