#trivy:ignore:AVD-AWS-0010 ignore warning about logging
#trivy:ignore:AVD-AWS-0011 ignore warning about waf
resource "aws_cloudfront_distribution" "ec2_distribution" {
  origin {
    domain_name = aws_instance.this.public_dns
    origin_id   = aws_instance.this.public_dns

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = ""

  aliases = [var.domain_name]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_instance.this.public_dns

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_All"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}

#trivy:ignore:AVD-AWS-0010 ignore warning about logging
#trivy:ignore:AVD-AWS-0011 ignore warning about waf
resource "aws_cloudfront_distribution" "webhook_distribution" {
  origin {
    domain_name = aws_instance.this.public_dns
    origin_id   = aws_instance.this.public_dns

    custom_origin_config {
      http_port              = var.webhook_port
      https_port             = var.webhook_port
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  comment             = "Webhook for notifying server of upgrade to app"
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = ""

  aliases = [var.domain_name]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_instance.this.public_dns

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_All"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}
