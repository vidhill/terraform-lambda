provider "aws" {
  region = var.region
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "resize/function.zip"
  function_name = "vidhill-resize-lambda"
  role          = "arn:aws:iam::728615433596:role/iam_for_lambda"
  handler       = "index.handler"
  runtime       = "nodejs12.x"
}

