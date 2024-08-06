#! /bin/bash

# Variables for credentials
DB_USER="telex_be_user"
DB_PASSWORD="p@$$w0rd"
DB_NAME="telex_be_db"

# Create the database and grant the user access to it
sudo -i -u postgres psql <<EOF
    -- Create a database named '$DB_NAME'
    CREATE DATABASE $DB_NAME;
    -- Create a user named '$DB_USER' with password '$DB_PASSWORD'
    CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD'
;
    -- Grant all privileges on the database '$DB_NAME' to the user '$DB_USER'
    GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

# Restart PostgreSQL to apply changes
sudo systemctl restart postgresql

echo "PostgreSQL setup is complete. User '$DB_USER' with database '$DB_NAME' has been created. The user can connect using the password '$DB_PASSWORD'."