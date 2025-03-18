# Set up Stately

## Step 1: Initial Setup

```sh
# Install the Stately CLI
curl -sL https://stately.cloud/install | sh

# Set up the demo repo
mkdir demo-w && cd demo-w
git init

# Create a schema
stately schema init ./schema-v1
```

## Step 2: Build a Schema

See `schema-v1/schema.ts`. This is a simple model for a permissions lease system.

* It leverages a single table design to store users, resources, and leases in a single store.
* Key paths are chosen to allow for:
    * Listing all active leases for a user. Users an see all their leases in the lease tool.
    * Listing all active leases for a resource. Useful for an admin panel or resource owner.
    * Getting the lease (or leases) for a specific user + resource. Useful for an AuthZ filter.
* A TTL on the lease means it expires a configurable time from when the lease was created.
* "touching" the lease extends the expiration.

```sh
# Crate a go.mod file - we need this to generate code
go mod init github.com/StatelyCloud/demo-w

# Generate preview code to see what it looks like
stately schema generate -l go --preview schema-v1/schema.ts pkg/schema

# Run mod tidy to "install" all the dependencies
go mod tidy
```

## Step 3: Develop a client

See the code in `pkg/client`. This implements some business logic around managing leases.

As a comparison we've also implemented a similar client in `pkg/ddb` that does the same thing directly in DDB. Note:

* It's much more code than the StatelyDB example. You need to write your own domain models and map them to DDB attributes.
* Single-table design is tricky to get right and the code is hard to understand, thanks to reuse of key names.
* This version doesn't have quite the same flexibility as the StatelyDB version - it uses GSIs to handle the alternate lookups, but doesn't put in place all the GSIs you might need.
* This version doesn't properly track createdAt and lastModifiedAt times.
* Validation needs to happen on the client side, since there's no schema to enforce shape.
* In the StatelyDB version, we easily enforce uniqueness of user by email - in the DDB version this requires carefully writing to (and reading from) two copies of the user with a transaction.

## Step 4: Set up a Store and Schema

1. Create the backing table with CloudFormation by following: https://docs.stately.cloud/deployment/byoc/:
   ```sh
   curl https://docs.stately.cloud/create-table.sh -o ./create-table.sh && chmod a+x ./create-table.sh
   ./create-table.sh demo-w
   ```
2. Create an IAM policy with restricted permissions:
   ```sh
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
   aws iam attach-role-policy \
      --role-name demo-w-role \
      --policy-arn arn:aws:iam::<ACCOUNT_ID>:policy/StatelyDBDynamoReadWriteAccess
   ```
3. In the Console, click "New Store", select BYOC and provide the table ARN from above.
4. Copy the store ID and schema ID. 6978411690784381 and 4291558376530788.
5. Create an Access Key with the type "Data Plane Key for BYOC" and copy the key string.
6. Put your schema:
   ```sh
   stately schema put -s 4291558376530788 schema-v1/schema.ts
   stately schema generate -l go -s 4291558376530788 pkg/schema
   ```

## Set up Kubernetes

1. Create an EKS cluster (offscreen)
    1. Create a service account for pod identity
    2. Associate it with our role:
       ```sh
       aws eks associate-pod-identity-profile \
          --cluster-name <CLUSTER_NAME> \
          --namespace default \
          --name demo-w-profile \
          --role-arn arn:aws:iam::<ACCOUNT_ID>:role/demo-w-role \
          --service-account demo-w
       ```
2. Create a k8s secret for the access key:
   ```sh
   kubectl create secret generic demo-w-secret --from-literal=STATELY_ACCESS_KEY="$STATELY_ACCESS_KEY"
   ```
3. Push your service container and create the deployment and loadbalancer using `k8s/deployment.yaml`