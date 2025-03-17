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
stately schema generate -l go --preview schema-v1/schema.ts pkg/schema/v1

# Run mod tidy to "install" all the dependencies
go mod tidy
```