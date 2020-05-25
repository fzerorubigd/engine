README
------

This is my code template for my web projects. I decided to make it open source, mostly for others to review.

Project uses a docker-compose based development environment. you need docker-compose, and by running `docker-compose up` (or `make start`) you can run the project. 

| Address | |
| ---| ---|
| localhost:8080/dashboard | Traefik dashboard |
|api.localhost| The API endpoint, For the API doc vist api.localhost/v1/swagger/ , Remember change the Scheme to HTTP for testing it locally |
|adminer.localhost| The adminer web interface, the database address is postgresql, user name and password are inside the`.env` file in root foldre|
|redisweb.localhost|The redis web UI|
|mailhog.localhost|The MailHog interface, not yet used in the project|

TODO: Need a real README!