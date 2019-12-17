README
------

This is my code template for my web projects. I decided to make it open source, mostly for others to review.

Currently there is a vagrant-docker configuration available with the project. the project requires postgresql/redis which is installed in the vagrant box. 

Run `vagrant up` and then `vagrant ssh`, then `cd engine`.

You can run `make test` for running the test and `make run-server-qollenge` then visit `localhost:8090/v1/swagger`

TODO: Need a real README!