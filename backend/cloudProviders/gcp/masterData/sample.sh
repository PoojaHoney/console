#!/bin/bash

#Root access
sudo su

# Variables
SSH_USER="admin"
SSH_PASSWORD="321@nimda#console-poc"
MONGO_USERNAME="lms-bo"
MONGO_DATABASE="lms-bo"
MONGO_PASSWORD="9wcEGYnFPgMWbsuQww0yTiEsQnW"
CLOUD_STORE_BUCKET="alph-a_console"
PROJECT_ID="peaceful-harbor-416116"
REGION="asia-south1"
PRODUCT="lms"
GCP_ARTIFACT_REGISTRY="$REGION-docker.pkg.dev/$PROJECT_ID/$PRODUCT"
SERVICE_ACCOUNT="lms-gcp-user"
FRONTEND_IMAGE="$GCP_ARTIFACT_REGISTRY/bo-frontend:0.0.2"
USER_IMAGE="$GCP_ARTIFACT_REGISTRY/bo-user:0.0.2"
CONTENT_IMAGE="$GCP_ARTIFACT_REGISTRY/bo-content:0.0.2"
INSTANCE_ID={{INSTANCE_ID}}

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



# Authenticate Docker with Artifact Registry
gcloud auth configure-docker $GCP_ARTIFACT_REGISTRY

gcloud auth print-access-token --impersonate-service-account $SERVICE_ACCOUNT@$PROJECT_ID.iam.gserviceaccount.com | docker login -u oauth2accesstoken --password-stdin https://$REGION-docker.pkg.dev

# Pull Docker images
docker pull $FRONTEND_IMAGE
docker pull $USER_IMAGE
docker pull $CONTENT_IMAGE

# Run MongoDB Docker container
docker run -d --name mongo -e MONGO_INITDB_ROOT_USERNAME=$MONGO_USERNAME -e MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASSWORD mongo

mkdir /home/$SSH_USER/$PRODUCT
mkdir /home/$SSH_USER/$PRODUCT/restore-data

# Import data into MongoDB
gsutil cp -r gs://$CLOUD_STORE_BUCKET/$PRODUCT/databases/mongo/backup/* /home/$SSH_USER/$PRODUCT/restore-data # Change path/on/your/local/machine
docker cp /home/$SSH_USER/$PRODUCT/restore-data mongo:/restore_data/
docker exec -i mongo mongorestore --host localhost:27017 --username=$MONGO_USERNAME --password=$MONGO_PASSWORD --authenticationDatabase=admin --db=$MONGO_DATABASE /restore_data

# Run Docker containers for services
docker run -d --name bo-frontend -p 4000:4000 $FRONTEND_IMAGE
docker run -d --name bo-user -p 4002:4002 $USER_IMAGE
docker run -d --name bo-content -p 4001:4001 $CONTENT_IMAGE