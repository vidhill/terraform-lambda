provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "srcBucket" {
  bucket = "vidhill-my-tf-test-bucket"

  tags = {
    Name = "Source bucket"
  }
}

resource "aws_s3_bucket" "destBucket" {
  bucket = "${aws_s3_bucket.srcBucket.bucket}-resized"
  tags = {
    Name = "Destination bucket"
  }
}

# Add permission for bucket to trigger lambda function
resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  principal     = "s3.amazonaws.com"
  function_name = aws_lambda_function.test_lambda.arn
  source_arn    = aws_s3_bucket.srcBucket.arn
}



locals {
  bucketIds = [
    aws_s3_bucket.srcBucket.id,
    aws_s3_bucket.destBucket.id
  ]
}

# Make buckets private
resource "aws_s3_bucket_acl" "example" {
  for_each = toset(local.bucketIds)
  bucket   = each.value
  acl      = "private"
}


# Deny public access to buckets
resource "aws_s3_bucket_public_access_block" "example" {
  for_each = toset(local.bucketIds)
  bucket   = each.value

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Add handler to trigger lambda on file add
resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.srcBucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.test_lambda.arn
    events              = ["s3:ObjectCreated:*"]
  }

}

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# Create policy allowing lambda read/write access to buckets: START
#

data "aws_iam_policy_document" "bucket_read_write" {

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
      "${aws_s3_bucket.srcBucket.arn}/*",
      "${aws_s3_bucket.destBucket.arn}/*"
    ]
  }
}

resource "aws_iam_role_policy" "resize_buckets_read_write_policy" {
  name   = "resize_buckets_read_write_policy"
  role   = aws_iam_role.iam_for_lambda.id
  policy = data.aws_iam_policy_document.bucket_read_write.json
}

#
# Create policy allowing lambda read/write access to buckets: END
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -



# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# Create role for lambda: START
#

data "aws_iam_policy_document" "assume_role" {
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

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.iam_for_lambda.id
}

#
# Create role for lambda: END
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -


# Create lambda
resource "aws_lambda_function" "test_lambda" {
  filename      = data.archive_file.lambda_zip_dir.output_path
  function_name = "vidhill-resize-lambda"
  # role          = aws_iam_role.iam_for_lambda.arn
  role    = "arn:aws:iam::728615433596:role/iam_for_lambda"
  handler = "index.handler"
  runtime = "nodejs12.x"

  source_code_hash = filebase64sha256(data.archive_file.lambda_zip_dir.output_path)
}

#
# Create zip archive of lambda folder
#
data "archive_file" "lambda_zip_dir" {
  type        = "zip"
  output_path = "function.zip"
  source_dir  = data.external.build.working_dir
}

#
# Build (npm install in this case)
#
data "external" "build" {
  program = ["bash", "-c", <<EOT
    npm ci --arch=x64 --platform=linux --target=10.15.0 >&2 && echo "{}" 
  EOT
  ]
  working_dir = "${path.module}/resize"
}

