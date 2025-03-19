// Define your StatelyDB schema in this file!
// Check out our documentation at https://stately.cloud.

import {
  durationSeconds,
  itemType,
  migrate,
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
    '/user-:user_id/res-:resource_id/lease-:id',
    '/res-:resource_id/lease-:id',
    '/lease-:id',
  ],
  // Automatically delete leases after the time in the duration field since they
  // were last updated.
  ttl: {
    source: 'fromLastModified',
    field: 'duration_seconds',
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
    resource_id: {
      type: ResourceID,
    },
    /** Allow the user to specify why they needed the lease. */
    reason: {
      type: string,
      required: false,
    },
    /** Who has approved this? The lease is not considered valid until approved by another person. */
    approver: {
      type: UserID,
      required: false,
    },
    /** How long is this lease for? This is measured from when the lease was last modified. */
    duration_seconds: {
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

export const AddApprover = migrate(1, "Add approver and make reason optional", (m) => {
  m.changeType('Lease', (t) => {
    t.addField('approver');
    t.renameField('res_id', 'resource_id');
    t.renameField('duration', 'duration_seconds');
  })
});

export const ReasonNotRequired = migrate(2, "Make reason not required", (m) => {
  m.changeType('Lease', (t) => {
    t.markFieldAsNotRequired('reason', 'No reason given');
  })
});