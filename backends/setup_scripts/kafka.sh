#!/usr/bin/env bash
## Start a test Kafka instance
## This is used for testing only, this doesn't build a prod ready Kafka instance.
## Make sure you change the below if you're actually using this!
SSLPW=changeme!

apt install default-jre
wget http://mirror.ventraip.net.au/apache/kafka/2.1.0/kafka_2.11-2.1.0.tgz
tar -xzf kafka_2.11-2.1.0.tgz
cd kafka_2.11-2.1.0
bin/zookeeper-server-start.sh config/zookeeper.properties
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
# holy shit I hate java look at how hard this is to trust ONE CA certificate.
#Step 1
keytool -keystore server.keystore.jks -alias localhost -validity 365 -keyalg RSA -genkey
#Step 2
openssl req -new -x509 -keyout ca-key -out ca-cert -days 365
keytool -keystore server.truststore.jks -alias CARoot -import -file ca-cert
keytool -keystore client.truststore.jks -alias CARoot -import -file ca-cert
#Step 3
keytool -keystore server.keystore.jks -alias localhost -certreq -file cert-file
openssl x509 -req -CA ca-cert -CAkey ca-key -in cert-file -out cert-signed -days 365 -CAcreateserial -passin pass:$SSLPW
keytool -keystore server.keystore.jks -alias CARoot -import -file ca-cert
keytool -keystore server.keystore.jks -alias localhost -import -file cert-signed
printf " !! If you're using a server with less than 1G memory\n";
printf " !! change line 29 in bin/kafka-start-server.sh to reflect a more reasonable amount\n";
printf " !!\n";
printf " !! Also, if you're running on AWS, change advertised.listeners=PLAINTEXT://[ your actual servers name ]]:9092 in server.config."
printf "!!\n";
printf "!! For SSL Support, add the following snippet:\n\n";
echo "listeners=PLAINTEXT://:9092, SSL://:9093
ssl.keystore.location=/home/ubuntu/kafka_2.11-2.1.0/server.keystore.jks
ssl.keystore.password=$SSLPW
ssl.key.password=$SSLPW
ssl.truststore.location=/home/ubuntu/kafka_2.11-2.1.0/server.truststore.jks
ssl.truststore.password=spaghett
"




