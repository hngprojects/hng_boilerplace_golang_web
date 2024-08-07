#!/bin/bash

set -e

export BRANCH="dev"
export APP_NAME="telex_be"
export APPROOT=~/deployments/telex_be
export PATH=$PATH:~/.nvm/versions/node/v20.15.1/bin
export PATH=$PATH:/usr/local/go/bin

mkdir -p $APPROOT
cd $APPROOT

if [ -d "$APPROOT/.git" ]; then
  # Navigate to the repository directory and pull changes
  cd $APPROOT || { echo "Failed to navigate to web root directory"; exit 1; }
  git reset --hard HEAD || { echo "Failed to reset local changes"; exit 1; }
  git pull origin $BRANCH || { echo "Failed to pull latest changes"; exit 1; }
else
  echo "Application not on server, consider cloning it"
  exit 1
fi

# Replace environment variables in app.env
cp app-sample.env app.env
for var in "$@"
do
    # Split the variable into key and value
    KEY=${var%%=*}
    VALUE=${var#*=}

    # Check if the key already exists in the file
    if grep -q "^$KEY=" app.env; then
        # Update the existing key with the new value
        sed -i "s|^$KEY=.*|$KEY=$VALUE|" app.env
    else
        # Add the new key-value pair
        echo -e "\n$KEY=$VALUE" >> app.env
    fi
done

go build -o $APPROOT/$APP_NAME

# Start or restart the built binary with pm2
if pm2 list | grep -qw "$APP_NAME"; then
  echo "Process $APP_NAME is running. Restarting..."
  pm2 restart $APPROOT/$APP_NAME
else
  echo "Process $APP_NAME is not running. Starting..."
  pm2 start $APPROOT/$APP_NAME --name $APP_NAME
fi
