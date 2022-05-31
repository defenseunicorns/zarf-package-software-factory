output "public_instance_id" {
  value = aws_spot_instance_request.public.id
}

output "public_instance_ip" {
  value = aws_spot_instance_request.public.public_ip
}
