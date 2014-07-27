[![Stories in Ready](https://badge.waffle.io/granthco/granth.co.png?label=ready&title=Ready)](https://waffle.io/granthco/granth.co)

[![Build Status](https://drone.io/github.com/GranthCo/granth.co/status.png)](https://drone.io/github.com/GranthCo/granth.co/latest)

Granth.co source code
======================
This is the source code for http://granth.co 


Instructions
============

Download datasource from: 

http://www.sikher.com/sql/2.x/


You can set the following environment variables to override the defaults.

  GRANTHCO_DATABASE_USERNAME="`username`"

  GRANTHCO_DATABASE_PASSWORD="`password`"

  GRANTHCO_DATABASE_HOST="`host`"

  GRANTHCO_DATABASE_PORT="`port`"

  GRANTHCO_DATABASE_NAME="`database`"

The default values used are *root*, *password*, *localhost*, *3306* and *gurbanidb* respectively

---

For compiling:

    GOPATH=/path/to/granth.co/Godeps/_workspace/ go build

And then run the app as:

    PORT=8888 ./granth.co
