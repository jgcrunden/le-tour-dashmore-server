terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 6.6.0"
    }

    http = {
      source = "hashicorp/http"
      version = ">= 3.5.0"
    }
  }
  required_version = ">=1.12.2"
  backend "s3" {
    bucket = "le-tour-terraform-state"
    key = "letour.tfstate"
    region = "eu-west-1"
  }
}

provider "aws" {
  alias  = "acm_provider"
  region = "us-east-1"
}
