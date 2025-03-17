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

## Step 3: Set up a Store and Schema

1. Create the backing table and IAM roles in your account by following: https://docs.stately.cloud/deployment/byoc/
2. In the Console, click "New Store", select BYOC and provide the table ARN from above.
3. Copy the store ID and schema ID. 6978411690784381 and 4291558376530788.
5. Create an Access Key with the type "Data Plane Key for BYOC".

## Step 4: Develop a client

See the code in `pkg/client`. This implements some business logic around managing leases.

As a comparison we've also implemented a similar client in `pkg/ddb` that does the same thing directly in DDB. Note:

* It's much more code than the StatelyDB example. You need to write your own domain models and map them to DDB attributes.
* Single-table design is tricky to get right and the code is hard to understand, thanks to reuse of key names.
* This version doesn't have quite the same flexibility as the StatelyDB version - it uses GSIs to handle the alternate lookups, but doesn't put in place all the GSIs you might need.
* This version doesn't properly track createdAt and lastModifiedAt times.
* Validation needs to happen on the client side, since there's no schema to enforce shape.
* In the StatelyDB version, we enforce uniqueness of user by email - in the DDB version this is very difficult to do.