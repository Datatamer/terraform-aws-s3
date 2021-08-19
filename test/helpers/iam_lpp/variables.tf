variable "name_prefix" {
  type    = string
  default = "test"
}

variable "policies_arn" {
  type = list(string)
}