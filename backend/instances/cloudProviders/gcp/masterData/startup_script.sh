#!/bin/bash

#Root access
sudo su

# Variables
SSH_USER="admin"
SSH_PASSWORD="321@nimda#console-poc"
MONGO_USERNAME="lms-bo"
MONGO_DATABASE="lmsbo"
MONGO_PASSWORD="9wcEGYnFPgMWbsuQww0yTiEsQnW"
CLOUD_STORE_BUCKET={{CLOUD_STORE_BUCKET}}
PROJECT_ID={{PROJECT_ID}}
REGION={{REGION}}
INSTANCE_ID={{INSTANCE_ID}}
INSTANCE_NAME={{INSTANCE_NAME}}
USER_IMAGE_VERSION={{USER_IMAGE_VERSION}}
CONTENT_IMAGE_VERSION={{CONTENT_IMAGE_VERSION}}
FRONTEND_IMAGE_VERSION={{FRONTEND_IMAGE_VERSION}}
USER_NAME={{USER_NAME}}
CONTENT_NAME={{CONTENT_NAME}}
FRONTEND_NAME={{FRONTEND_NAME}}
USER_PORT={{USER_PORT}}
CONTENT_PORT={{CONTENT_PORT}}
FRONTEND_PORT={{FRONTEND_PORT}}
PRODUCT="lms"
GCP_ARTIFACT_REGISTRY="$REGION-docker.pkg.dev/$PROJECT_ID/$PRODUCT"
SERVICE_ACCOUNT={{SERVICE_ACCOUNT}}
FRONTEND_IMAGE="$GCP_ARTIFACT_REGISTRY/$FRONTEND_NAME:$FRONTEND_IMAGE_VERSION"
USER_IMAGE="$GCP_ARTIFACT_REGISTRY/$USER_NAME:$USER_IMAGE_VERSION"
CONTENT_IMAGE="$GCP_ARTIFACT_REGISTRY/$CONTENT_NAME:$CONTENT_IMAGE_VERSION"
export PUBLIC_IP=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip)

# Add a new user to the VM with sudo privileges
useradd -m -s /bin/bash $SSH_USER
echo "$SSH_USER:$SSH_PASSWORD" | chpasswd
usermod -aG sudo $SSH_USER

# Install Docker
apt-get update
apt-get install -y docker.io

#SSH Keys Generation
su - $SSH_USER -c 'ssh-keygen -t rsa -b 2048 -f ~/.ssh/id_rsa -N ""'

# SSH Key Authentication
su - $SSH_USER -c 'cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys'
su - $SSH_USER -c 'chmod 0600 ~/.ssh/authorized_keys'

gsutil cp /home/$SSH_USER/.ssh/id_rsa gs://$CLOUD_STORE_BUCKET/ssh-keys/

sudo systemctl restart sshd

gsutil cp gs://$CLOUD_STORE_BUCKET/$PRODUCT/$INSTANCE_ID/serviceAccounts/keys/$INSTANCE_ID$INSTANCE_NAME.json key.json
gcloud auth activate-service-account --key-file key.json

# Authenticate Docker with Artifact Registry
gcloud auth configure-docker $GCP_ARTIFACT_REGISTRY

gcloud auth print-access-token --impersonate-service-account $SERVICE_ACCOUNT | docker login -u oauth2accesstoken --password-stdin https://$REGION-docker.pkg.dev

# Pull Docker images
docker pull $FRONTEND_IMAGE
docker pull $USER_IMAGE
docker pull $CONTENT_IMAGE
docker pull mongo:latest

docker network create $PRODUCT

# Run MongoDB Docker container
docker run -d --name mongo -e MONGO_INITDB_ROOT_USERNAME=$MONGO_USERNAME -e MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASSWORD --network=$PRODUCT mongo

mkdir /home/$SSH_USER/$PRODUCT
mkdir /home/$SSH_USER/$PRODUCT/restore-data

# Import data into MongoDB
gsutil cp -r gs://$CLOUD_STORE_BUCKET/$PRODUCT/databases/mongo/backup/* /home/$SSH_USER/$PRODUCT/restore-data # Change path/on/your/local/machine
docker cp /home/$SSH_USER/$PRODUCT/restore-data mongo:/restore_data/
docker exec -i mongo mongorestore --host localhost:27017 --username=$MONGO_USERNAME --password=$MONGO_PASSWORD --authenticationDatabase=admin --db=$MONGO_DATABASE /restore_data

# Run Docker containers for services
docker run -d --name $FRONTEND_NAME -p $FRONTEND_PORT:$FRONTEND_PORT --network=$PRODUCT -e CONTENT_API_HOST=$PUBLIC_IP -e USERS_API_HOST=$PUBLIC_IP $FRONTEND_IMAGE
docker run -d --name $USER_NAME -p $USER_PORT:$USER_PORT --network=$PRODUCT $USER_IMAGE
docker run -d --name $CONTENT_NAME -p $CONTENT_PORT:$CONTENT_PORT --network=$PRODUCT $CONTENT_IMAGE
