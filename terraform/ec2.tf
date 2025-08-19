data "aws_ami" "this" {
  most_recent = true
  //  owners      = ["amazon"]
  owners = ["125523088429"]
  filter {
    name   = "architecture"
    values = ["arm64"]
  }
  filter {
    name   = "name"
    values = ["Fedora-Cloud-Base-AmazonEC2.aarch64-42-1.1"]
  }
}

resource "aws_instance" "this" {
  ami = data.aws_ami.this.id
  tags = {
    Name = "le-tour-dashmore-server"
  }
  instance_market_options {
    market_type = "spot"
  }
  instance_type = "t4g.micro"
  metadata_options {
    http_tokens = "required"
  }
  root_block_device {
    encrypted = true
  }
  key_name               = aws_key_pair.ssh_access.id
  vpc_security_group_ids = [aws_security_group.this.id, aws_security_group.webhook.id]

  user_data = file("user_data.sh")
}

resource "aws_key_pair" "ssh_access" {
  public_key = var.public_key
  key_name   = "ssh-key"
}

resource "aws_security_group" "this" {
  name        = "Security group for EC2 instance"
  description = "Allow SSH inbound traffic, http inbound from cloudfront and all outbound traffic"
}

resource "aws_security_group" "webhook" {
  name        = "Security group for Webhook to EC2 instance"
  description = "Allow webhook acces from cloudfront"
}

data "http" "myip" {
  url = "https://ipv4.icanhazip.com"
}

locals {
  local_public_ip = "${chomp(data.http.myip.response_body)}/32"
}

data "aws_ec2_managed_prefix_list" "cloudfront" {
  filter {
    name   = "owner-id"
    values = ["AWS"]
  }

  filter {
    name   = "prefix-list-name"
    values = ["com.amazonaws.global.cloudfront.origin-facing"]
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh_ipv4" {
  description       = "Allow ssh access to ec2 instance"
  security_group_id = aws_security_group.this.id
  cidr_ipv4         = local.local_public_ip
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "allow_http_ipv4" {
  description       = "Application Allow http access to ec2 instance from cloudfront only"
  security_group_id = aws_security_group.this.id

  prefix_list_id = data.aws_ec2_managed_prefix_list.cloudfront.id
  from_port      = 80
  ip_protocol    = "tcp"
  to_port        = 80
}

resource "aws_vpc_security_group_ingress_rule" "webhook_allow_http_ipv4" {
  description       = "Webhook - Allow http access to ec2 instance from cloudfront only"
  security_group_id = aws_security_group.webhook.id

  prefix_list_id = data.aws_ec2_managed_prefix_list.cloudfront.id
  from_port      = var.webhook_port
  ip_protocol    = "tcp"
  to_port        = var.webhook_port
}

#trivy:ignore:AVD-AWS-0104 ignore warning about egress
resource "aws_vpc_security_group_egress_rule" "egress" {
  description       = "Egress from EC2 instance"
  security_group_id = aws_security_group.this.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = -1
  ip_protocol       = "-1"
  to_port           = -1
}

output "ssh_command" {
  value = "ssh fedora@${aws_instance.this.public_ip}"
}
