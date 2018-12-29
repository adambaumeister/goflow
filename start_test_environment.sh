#!/usr/bin/env bash
# Start the venv we will install awscli into
# This is really just to make sure the process of building a test env works OK
sudo apt-get install python3-venv
python3 -m venv aws-cli-tools
cd aws-cli-tools
. bin/activate
pip3 install awscli
# The following env variables are stored in Travis:
#   - AWS_ACCESS_KEY_ID         : AWS access key
#   - AWS_SECRET_ACCESS_KEY     : AWS Secret
#   - AWS_INSTANCE_ID           : Instance to test against
#   - SQL_PASSWORD              : The SQL Database password for backend testing
export AWS_DEFAULT_REGION=ap-southeast-2
aws ec2 start-instances --instance-ids $AWS_INSTANCE_ID
# Wait for it to start
sleep 90
export SQL_SERVER=`aws ec2 describe-instances --instance-ids $AWS_INSTANCE_ID | grep PublicDnsName | head -n 1 | awk -F\" '{print $4}'`
nc -w2 -z $SQL_SERVER 5432
if [ $? -ne 0 ]
then
    echo "Failed to start instance: $AWS_INSTANCE_ID"
    exit 1
fi
echo $SQL_SERVER
cd ..
