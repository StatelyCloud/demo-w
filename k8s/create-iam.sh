#!/usr/bin/env bash

aws iam create-policy --policy-name StatelyDBDynamoReadWriteAccess \
  --policy-document '{
        "Version": "2012-10-17",
        "Statement": [{
          "Effect": "Allow",
          "Action": [
            "dynamodb:*Item",
            "dynamodb:Describe*",
            "dynamodb:List*",
            "dynamodb:Query",
            "dynamodb:Scan"
          ],
          "Resource" : "*"
        }]}'
POLICY_ARN=$(aws iam list-policies --query 'Policies[?PolicyName==`StatelyDBDynamoReadWriteAccess`].Arn' --output text)
aws iam create-role \
  --role-name demo-w-role \
  --assume-role-policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "pods.eks.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }'
ROLE_ARN=$(aws iam get-role --role-name demo-w-role --query 'Role.Arn' --output text)

# Set up pod identity
aws iam attach-role-policy \
  --role-name demo-w-role \
  --policy-arn $POLICY_ARN
aws eks associate-pod-identity-profile \
  --cluster-name $CLUSTER_NAME \
  --namespace default \
  --name demo-w-profile \
  --role-arn $ROLE_ARN \
  --service-account demo-w
