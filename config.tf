provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "bucket"{
  bucket = "vidhill-my-tf-test-bucket"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket" "bucket1"{
  bucket =  "${aws_s3_bucket.bucket.bucket}-resized"
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
      "arn:aws:s3:::vidhill-my-tf-test-bucket/*",
      "arn:aws:s3:::vidhill-my-tf-test-bucket-resized/*"
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
  name = "test_policy"
  role = aws_iam_role.iam_for_lambda.id

  policy = data.aws_iam_policy_document.example1.json
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = data.aws_iam_policy_document.example2.json
}

resource "aws_lambda_function" "test_lambda" {
  filename      = "resize/function.zip"
  function_name = "vidhill-resize-lambda"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.handler"
  runtime       = "nodejs12.x"
}

