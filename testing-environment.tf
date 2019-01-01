/*
testing-environment

Provides an environment for testing potential goflow releases.

Builds a server and vpc config but at this time does not run post-provisioning (setup_scripts).
Requires:
 env vars:
  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_KEY

Outputs:
  - Public DNS endpoint
  - Instance ID
*/
provider "aws" {
  region = "ap-southeast-2"
}

terraform {
  backend "local" {
    path = "..\\..\\..\\..\\..\\..\\..\\Temp\\terraform.tfstate"
  }
}

variable "az" {
  default = "ap-southeast-2a"
}

variable "vpc_config" {
  type = "map",
  default = {
    network_block = "10.10.0.0/16",
    name = "goflow-test-vpc",
  }
}

variable "trusted-ips" {
  type = "list"
  default = ["124.171.205.1/32"]
}

variable "public_network_config" {
  type = "map",
  default = {
    network = "10.10.10.0/24",
    name = "goflow-test-subnet",
  }
}

data "aws_route53_zone" "selected" {
  name         = "spaghettsucks.com."
}

resource "aws_route53_record" "goflow-test" {
  name = "goflow-test.${data.aws_route53_zone.selected.name}"
  type = "CNAME"
  ttl = "300"
  zone_id = "${data.aws_route53_zone.selected.zone_id}"
  records = ["${aws_instance.goflow-test-instance.public_dns}"]
}

resource "aws_key_pair" "kp" {
  key_name = "home-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAk9wAsRVSw0dWgcukpiamRe9KwWe5Itr1mKMXWcQZw2M6XwY9w67geUavuLOeia/GgcUax/fWa8Z2hwL3S8q2TuyBvpaWvKLfSWPm+U0zsNE2rFGa2Iw/rQxZUEL1N8njoPN+LURr47iZNATbtSRIDyUlpGF5ejUsDHRpvGD4VJTgYXP6scdgazqAFsex9pElvu1dLHFkzAOuHYEnFTkcPHTSTtgujoz82iFa0GfpTQ0wcI63/8dFRHKN+BEgWgdx35xHzQ9QLXwvVnEVWs5T6bg5Wvjd5sF6AdCcdfiQuzZpCGr/o6KjJAoE8yvBzzo1ycN8ulo8THlMshzhSqrn1Q== rsa-key-20180112"
}

resource "aws_instance" "goflow-test-instance" {
  ami = "ami-07a3bd4944eb120a0"
  instance_type = "t2.micro"
  vpc_security_group_ids = ["${aws_vpc.main_vpc.default_security_group_id}"]
  key_name = "${aws_key_pair.kp.key_name}"
  subnet_id = "${aws_subnet.public_subnet.id}"

  tags {
    Name = "goflow-test-pgsql-instance"
  }
}

resource "aws_vpc" "main_vpc" {
  cidr_block  = "${var.vpc_config["network_block"]}"
  enable_dns_hostnames = true
  tags {
    Name = "${var.vpc_config["name"]}"
  }
}

resource "aws_subnet" "public_subnet" {
  vpc_id = "${aws_vpc.main_vpc.id}"
  cidr_block = "${var.public_network_config["network"]}"
  map_public_ip_on_launch = true
  availability_zone = "${var.az}"
  tags {
    Name = "${var.public_network_config["name"]}"
  }
}

/*
Internet Gateway
Provides inbound and outbound access to the internet
*/
resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.main_vpc.id}"

  tags {
    Name = "${var.vpc_config["name"]}_igw"
  }
}

resource "aws_route" "internet_access" {
  route_table_id         = "${aws_vpc.main_vpc.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.gw.id}"
}

/*
Security group rule
Configure a rule to allow access within a security group
*/
resource "aws_security_group_rule" "allow_all_from_trusted" {
  type            = "ingress"
  from_port       = 0
  to_port         = 65535
  protocol        = "tcp"
  cidr_blocks     = "${var.trusted-ips}"
  security_group_id = "${aws_vpc.main_vpc.default_security_group_id}"
}
resource "aws_security_group_rule" "allow_pgsql_from_any" {
  type            = "ingress"
  from_port       = 5432
  to_port         = 5432
  protocol        = "tcp"
  cidr_blocks     = ["0.0.0.0/0"]
  security_group_id = "${aws_vpc.main_vpc.default_security_group_id}"
}

resource "aws_security_group_rule" "allow_mysql_from_any" {
  type            = "ingress"
  from_port       = 3306
  to_port         = 3306
  protocol        = "tcp"
  cidr_blocks     = ["0.0.0.0/0"]
  security_group_id = "${aws_vpc.main_vpc.default_security_group_id}"
}

output "instance-dns" {
  value = "${aws_instance.goflow-test-instance.public_dns}"
}
output "instance-id" {
  value = "${aws_instance.goflow-test-instance.id}"
}