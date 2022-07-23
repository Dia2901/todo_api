# todo_api

Framework for API: Gin
Database : postgres sql
Cache : redis

_________________________________________________________________________________________________________________________

To Connect with the postgres database

1. Install postgres sql
2. Start postgres sql server
3. Create a database names postgres
4. Create table using
        CREATE TABLE tasks (
            id SERIAL PRIMARY KEY,
            des TEXT,
            status bool,
            username TEXT
        )

_________________________________________________________________________________________________________________________

To connect with redis 

1.Install redis
2. Run "redis-server"
3.In a new terminal window run "redis-cli"
4. Run "Select 1"

_________________________________________________________________________________________________________________________

To run the API

1. From the command line in the directory containing main.go, run the code using "go run ."
2. From a new command line window, use curl to make a request to your running web service.
    See Requests.txt
 
_________________________________________________________________________________________________________________________ 


