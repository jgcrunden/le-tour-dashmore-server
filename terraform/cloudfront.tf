#trivy:ignore:AVD-AWS-0010 ignore warning about logging
#trivy:ignore:AVD-AWS-0011 ignore warning about waf
resource "aws_cloudfront_distribution" "ec2_distribution" {
  origin {
    domain_name = aws_instance.this.public_dns
    origin_id   = aws_instance.this.public_dns

    custom_origin_config {
      http_port              = 8080
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
    default_ttl            = 0
    max_ttl                = 0
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
      http_port              = 9000
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  comment             = "Webhook for notifying server of upgrade to app"
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = ""

  aliases = [format("webhook.%s", var.domain_name)]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = aws_instance.this.public_dns

    viewer_protocol_policy   = "redirect-to-https"
    min_ttl                  = 0
    default_ttl              = 0
    max_ttl                  = 0
    cache_policy_id          = aws_cloudfront_cache_policy.query_strings.id
    origin_request_policy_id = aws_cloudfront_origin_request_policy.query_strings.id
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

resource "aws_cloudfront_origin_request_policy" "query_strings" {
  name = "le-tour-webhook-query-string-policy"
  headers_config {
    header_behavior = "none"
  }

  query_strings_config {
    query_string_behavior = "whitelist"
    query_strings {
      items = ["token"]
    }
  }

  cookies_config {
    cookie_behavior = "none"
  }
}

resource "aws_cloudfront_cache_policy" "query_strings" {
  name        = "le-tour-webhook-cache-policy"
  default_ttl = 0
  max_ttl     = 0
  min_ttl     = 0
  parameters_in_cache_key_and_forwarded_to_origin {
    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "none"
    }

    cookies_config {
      cookie_behavior = "none"
    }
  }

}

