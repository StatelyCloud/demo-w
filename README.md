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
            "Action": [
                "sts:AssumeRole",
                "sts:TagSession"
            ]
        }
        ]
      }'
   aws iam attach-role-policy \
      --role-name demo-w-role \
      --policy-arn arn:aws:iam::<ACCOUNT_ID>:policy/StatelyDBDynamoReadWriteAccess
   ```
3. In the Console, click "New Store", select BYOC and provide the table ARN from above.
4. Copy the store ID and schema ID:
   ```sh
   export STATELY_STORE_ID=4811130409281414
   export STATELY_SCHEMA_ID=4291558376530788
   ```
5. Create an Access Key with the type "Data Plane Key for BYOC" and copy the key string.
6. Put your schema:
   ```sh
   stately schema put -s $STATELY_SCHEMA_ID schema-v1/schema.ts
   stately schema generate -l go -v 1 -s $STATELY_SCHEMA_ID pkg/schema
   ```

## Step 5: Set up Kubernetes

1. Create an EKS cluster (offscreen)
    1. Create a service account for pod identity:
       ```sh
       kubectl apply -f k8s/service-account.yaml
       ```
    2. Associate it with our role:
       ```sh
       aws eks create-pod-identity-association \
          --cluster-name $CLUSTER_NAME \
          --namespace default \
          --role-arn arn:aws:iam::<ACCOUNT_ID>:role/demo-w-role \
          --service-account demo-w
       ```
2. Create a k8s secret for the access key:
   ```sh
   kubectl create secret generic demo-w-secret --from-literal=STATELY_ACCESS_KEY="$STATELY_ACCESS_KEY"
   ```
3. Push your service container and create the deployment and loadbalancer using `k8s/deployment.yaml`:
   ```sh
   ./publish.sh
   kubectl apply -f k8s/deployment.yaml
   ```

## Step 6: Try out our service

Test the service with these curl commands:

```sh
export DEMO_HOST="ac049c19b626845c9a3f9cc15ae94220-2137334291.us-west-2.elb.amazonaws.com"

# Create a user
curl -X POST http://$DEMO_HOST/users \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com", "name":"John Doe"}'

# Create a resource
curl -X POST http://$DEMO_HOST/resources \
  -H "Content-Type: application/json" \
  -d '{"name":"sensitive-database"}'

# Create a lease (replace UUIDs with actual IDs from previous responses)
curl -X POST http://$DEMO_HOST/leases \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "FY4wCvQLT9ycXM0jmv3nTg==",
    "resourceId": "uBrp9ZP8SR6WvcKYL8WCLg==",
    "reason": "Database maintenance",
    "durationHours": 0.5
  }'

# Get leases for a user (replace UUID with actual user ID)
curl http://$DEMO_HOST/users/FY4wCvQLT9ycXM0jmv3nTg==

# Get leases for a resource (replace UUID with actual resource ID)
curl http://$DEMO_HOST/resources/uBrp9ZP8SR6WvcKYL8WCLg==
```

Replace `localhost:8080` with your actual service URL if deploying to Kubernetes.

## Step 7: Updating schema

1. In `schema-v2/schema.ts` we've renamed some fields and added an approver.
2. Publish a new version of the schema:
   ```sh
   stately schema put -s $STATELY_SCHEMA_ID schema-v2/schema.ts
   stately schema generate -l go -v 2 -s $STATELY_SCHEMA_ID pkg/schema
   ```
3. Validate that the original service is still returning the old shapes
4. Publish a new version of the service
5. Show that the cURL commands return the new shape
6. TODO: how to keep the old version around? Maybe just have a separate deployment?
