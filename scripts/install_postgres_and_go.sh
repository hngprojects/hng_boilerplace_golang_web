#! /bin/bash

# Update package lists
sudo apt-get update

# Install GoLang
sudo snap install go --classic

# Install PostgreSQL
sudo apt-get install -y postgresql postgresql-contrib

# Start & Enale PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql