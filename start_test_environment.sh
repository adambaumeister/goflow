#!/usr/bin/env bash
#
# Start_test_environment
# Using the AWS cli, starts a preconfigured EC2 instance.
#
# Requires:
# ENV VARS
#   - AWS_ACCESS_KEY_ID         : AWS access key ; sourced from AWS (manual config)
#   - AWS_SECRET_ACCESS_KEY     : AWS Secret ; sourced from AWS (manual config)
#   - AWS_INSTANCE_ID           : Instance to test against ; Sourced from AWS ; see testing-environment.tf
#   - ZONE_ID                   : Hosted DNS Zone id; Sourced from AWS ; see testing-environment.tf
#   - SQL_PASSWORD              : The SQL Database password for backend testing ; Sourced from setup scripts.
# Start the venv we will install awscli into
# This is really just to make sure the process of building a test env works OK
python3 -m venv aws-cli-tools
cd aws-cli-tools
. bin/activate
pip3 install awscli
export AWS_DEFAULT_REGION=ap-southeast-2
aws ec2 start-instances --instance-ids $AWS_INSTANCE_ID
# Wait for it to start
sleep 30
# Update the dns record
aws route53 change-resource-record-sets --hosted-zone-id $ZONE_ID --change-batch file://rr_update.json
# Wait 60 seconds for that to finish...
sleep 60
export SQL_SERVER=goflow-test.spaghettsucks.com
nc -w2 -z $SQL_SERVER 5432
if [ $? -ne 0 ]
then
    echo "Failed to start instance: $AWS_INSTANCE_ID"
    exit 1
fi
echo $SQL_SERVER
cd ..
