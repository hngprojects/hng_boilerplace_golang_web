#! /bin/bash

# Variables for credentials
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="db_name"

# Check if the DB_USER is postgres and alter the password of this user else create a new user
if [ "$DB_USER" == "postgres" ]; then
  sudo -i -u postgres psql <<EOF
    ALTER USER postgres WITH PASSWORD '$DB_PASSWORD';
EOF
    echo "Password for 'postgres' user has been updated to '$DB_PASSWORD'."
else
    sudo -i -u postgres psql <<EOF
    -- Create a user named '$DB_USER' with password '$DB_PASSWORD'
    CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
EOF
fi

# Create the database and grant the user access to it
sudo -i -u postgres psql <<EOF 
    -- Create a database named '$DB_NAME'
    CREATE DATABASE $DB_NAME;

    -- Grant all privileges on the database '$DB_NAME' to the user '$DB_USER'
    GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

# Modify pg_hba.conf to allow password authentication
PG_HBA_FILE=$(sudo -i -u postgres psql -t -P format=unaligned -c 'SHOW hba_file')
sudo sed -i "s/^local\s\+all\s\+all\s\+peer/local   all             all                                     md5/" $PG_HBA_FILE
sudo sed -i "s/^host\s\+all\s\+all\s\+127.0.0.1\/32\s\+ident/host    all             all             127.0.0.1\/32            md5/" $PG_HBA_FILE
sudo sed -i "s/^host\s\+all\s\+all\s\+::1\/128\s\+ident/host    all             all             ::1\/128                 md5/" $PG_HBA_FILE
sudo bash -c "echo 'host    all             all             0.0.0.0/0               md5' >> $PG_HBA_FILE"
sudo bash -c "echo 'host    all             all             ::/0                    md5' >> $PG_HBA_FILE"

# Modify postgresql.conf to listen on all addresses
POSTGRESQL_CONF=$(sudo -i -u postgres psql -t -P format=unaligned -c 'SHOW config_file')
sudo sed -i "s/^#listen_addresses = 'localhost'/listen_addresses = '*'/" $POSTGRESQL_CONF

# Restart PostgreSQL to apply changes
sudo systemctl restart postgresql

echo "PostgreSQL setup is complete. User '$DB_USER' with database '$DB_NAME' has been created. The user can connect using the password '$DB_PASSWORD'."