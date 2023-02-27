provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "bucket" {
  bucket = "vidhill-my-tf-test-bucket"

  tags = {
    Name = "Source bucket"
  }
}

resource "aws_s3_bucket" "bucket1" {
  bucket = "${aws_s3_bucket.bucket.bucket}-resized"
  
  tags = {
    Name = "Destination bucket"
  }
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  principal     = "s3.amazonaws.com"
  function_name = aws_lambda_function.test_lambda.arn
  source_arn    = aws_s3_bucket.bucket.arn
}

locals {
  bucketIds = [
    aws_s3_bucket.bucket.id,
    aws_s3_bucket.bucket1.id
  ]
}

resource "aws_s3_bucket_acl" "example" {
  for_each = toset(local.bucketIds)
  bucket   = each.value
  acl      = "private"
}

resource "aws_s3_bucket_public_access_block" "example" {
  for_each = toset(local.bucketIds)
  bucket   = each.value

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.test_lambda.arn
    events              = ["s3:ObjectCreated:*"]
  }

}

data "aws_iam_policy_document" "example1" {

  statement {
    sid = ""
    actions = [
      "ec2:Describe*",
    ]
    effect    = "Allow"
    resources = ["*"]
  }

  statement {
    sid    = ""
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.bucket.arn}/*",
      "${aws_s3_bucket.bucket1.arn}/*"
    ]
  }
}

data "aws_iam_policy_document" "example2" {
  statement {
    sid = ""
    actions = [
      "sts:AssumeRole",
    ]
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com"
      ]
    }
  }
}


resource "aws_iam_role_policy" "test_policy" {
  name   = "test_policy"
  role   = aws_iam_role.iam_for_lambda.id
  policy = data.aws_iam_policy_document.example1.json
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.example2.json
}

resource "aws_iam_role_policy_attachment" "basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.iam_for_lambda.id
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "resize/function.zip"
  function_name = "vidhill-resize-lambda"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs12.x"
}

