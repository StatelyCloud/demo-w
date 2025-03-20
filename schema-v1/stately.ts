// Define your StatelyDB schema in this file!
// Check out our documentation at https://stately.cloud.

import {
  durationSeconds,
  itemType,
  string,
  timestampMilliseconds,
  type,
  uuid,
} from '@stately-cloud/schema';

// These are optional but help document ID types
export const UserID = type('UserID', uuid);
export const ResourceID = type('ResourceID', uuid);
export const LeaseID = type('LeaseID', uuid);

/**
 * A basic User object
 */
export const User = itemType('User', {
  keyPath: [
    '/user-:id',
    '/user_email-:email',
  ],
  fields: {
    id: {
      type: UserID,
      initialValue: 'uuid',
    },
    displayName: {
      type: string,
    },
    email: {
      type: string,
      valid: 'this.matches("[^@]+@[^@]+")',
    },
    createdAt: {
      type: timestampMilliseconds,
      fromMetadata: 'createdAtTime',
    },
  },
});

/**
 * A system is a resource that users can access.
 */
export const Resource = itemType('Resource', {
  keyPath: '/res-:id',
  fields: {
    id: {
      type: ResourceID,
      initialValue: 'uuid',
    },
    name: {
      type: string,
    },
    createdAt: {
      type: timestampMilliseconds,
      fromMetadata: 'createdAtTime',
    },
  },
});

/**
 * A "lease" gives users temporary access to a resource.
 */
export const Lease = itemType('Lease', {
  keyPath: [
    '/user-:user_id/res-:res_id/lease-:id',
    '/res-:res_id/lease-:id',
    '/lease-:id',
  ],
  // Automatically delete leases after the time in the duration field since they
  // were last updated.
  ttl: {
    source: 'fromLastModified',
    field: 'duration',
  },
  fields: {
    /** A unique identifier for the lease itself. */
    id: {
      type: LeaseID,
      initialValue: 'uuid',
    },
    /** The user that this lease is granted to. */
    user_id: {
      type: UserID,
    },
    /** The resource this lease grants access to. */
    res_id: {
      type: ResourceID,
    },
    /** Allow the user to specify why they needed the lease. */
    reason: {
      type: string,
    },
    /** How long is this lease for? This is measured from when the lease was last modified. */
    duration: {
      type: durationSeconds,
      required: false, // TODO: I would like this to be required though
    },
    /** Last touch time allows us to extend a lease by updating it. */
    lastTouched: {
      type: timestampMilliseconds,
      fromMetadata: 'lastModifiedAtTime',
    },
    createdAt: {
      type: timestampMilliseconds,
      fromMetadata: 'createdAtTime',
    },
  },
});
