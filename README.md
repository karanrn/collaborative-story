# Collaborative story

## How to run?
**Prerequisite:** - <br>
Run the **table_DDLs.sql** in the database for application code to run correctly.


If you are familiar (and fan) of docker (and docker-compose) - 
1. Edit the environment values in `.env` file
2. Run `docker-compose up -d` in the same directory as **docker-comose.yml** file
3. Now you can hit the API at `http://localhost:9000/`

And, if you want just run the code on your local machine - 
1. Export environment variable `VERLOOP_DSN`
    ```
    export VERLOOP_DSN="mysql://[user]:[password]@tcp([db_host]:3306)/[db_name]"
    ```
2. Now you can hit the API at `http://localhost:9000/`

**NOTE:**<br>
Application uses MySql as it's database, please pass the database connection string similar to below - <br>
`mysql://user:password@tcp(localhost:3306)/test`

---
## With more time...
I would 
1. Optimize structure and relations between the tables for least number of interactions from application (too many interactions with tables).
2. Smaller code and much more modularized (it is bit chunky now).
3. Better transaction handling in the application code.
4. Middlewares for authentication, rate limiting and request validation.
5. More and informal logging for the application.
6. Detailed profiling of the application.