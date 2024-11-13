# Gator
Blog and Website aggregator in go....

# Set Up
This requires latest golang installed.
Requires postresql installed. 
Create a ".gatorconfig.json" file in home directory 
run these commands 
   cd ~ 
   touch .gatorconfig.json
then edit with your editor, if you use vs code run
   code .gatorconfig.json 
save this json and edit to include your postgres db url
    "{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
 
    }"
make sure you create a database on your postgress named gator to do this run
    psql postgres://username:password@localhost:5432
Then run in psql command
    CREATE DATABASE gator;



# Installation 
Install gator with the command
    go install github.com/muhammadolammi/gator

# Available Commands 
 1. gator users
    This will list all users indicating the logged in user.
 2. gator register <user>
    This will register and auto login a  new user
 3. gator login <user>
    This will login to the provided user
 4. gator agg <time_between_reqs >
    run the aggregator that fetch feed every time_between_reqs 
 5. gator help
    list all commands and functionality

