resource "aws_route53domains_registered_domain" "this" {
  domain_name = var.domain_name

  dynamic "name_server" {
    for_each = toset(aws_route53_zone.le_tour_hosted_zone.name_servers)
    content {
      name = name_server.value

    }
  }
}

resource "aws_route53_zone" "le_tour_hosted_zone" {
  name = var.domain_name
}

resource "aws_route53_record" "certification_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = aws_route53_zone.le_tour_hosted_zone.zone_id
}

resource "aws_route53_record" "url_ip6" {
  name    = var.domain_name
  zone_id = aws_route53_zone.le_tour_hosted_zone.zone_id
  type    = "AAAA"

  alias {
    name                   = aws_cloudfront_distribution.ec2_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.ec2_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "www_ip4" {
  name    = var.domain_name
  zone_id = aws_route53_zone.le_tour_hosted_zone.zone_id
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.ec2_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.ec2_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "webhook_url_ip6" {
  name    = var.domain_name
  zone_id = aws_route53_zone.le_tour_hosted_zone.zone_id
  type    = "AAAA"

  alias {
    name                   = aws_cloudfront_distribution.webhook_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.webhook_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "webhook_www_ip4" {
  name    = var.domain_name
  zone_id = aws_route53_zone.le_tour_hosted_zone.zone_id
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.webhook_distribution.domain_name
    zone_id                = aws_cloudfront_distribution.webhook_distribution.hosted_zone_id
    evaluate_target_health = false
  }
}

