# Collaborative story

## How to run?
**Prerequisite:** - <br>
Run the **table_DDLs.sql** in the database for application code to run correctly.


If you are familiar (and fan) of docker (and docker-compose) - 
1. Edit the environment values in `.env` file
2. Run `docker-compose up -d` in the same directory as **docker-comose.yml** file
3. Now you can hit the API at `http://localhost:9000/`

And, if you want just run the code on your local machine - 
1. Export environment variable `DATABASE_DSN`
    ```
    export DATABASE_DSN="mysql://[user]:[password]@tcp([db_host]:3306)/[db_name]"
    ```
2. Run `go build -v -mod mod -ldflags "-s -w" -o restapi . `
3. Execute/Run `./restapi`
4. Now you can hit the API at `http://localhost:9000/`

**NOTE:**<br>
Application uses MySql as it's database, please pass the database connection string similar to below - <br>
`mysql://user:password@tcp(localhost:3306)/test`

---
## Things to be done...
I would 
1. Optimize structure and relations between the tables for least number of interactions from application (too many interactions with tables).
2. Middlewares for authentication, rate limiting and request validation.
3. More and informal logging for the application.
4. Detailed profiling of the application.

---
## Further reading - 
1. https://www.alexedwards.net/blog/organising-database-access
2. https://blog.lamida.org/mocking-in-golang-using-testify/
3. https://tutorialedge.net/golang/improving-your-tests-with-testify-go/