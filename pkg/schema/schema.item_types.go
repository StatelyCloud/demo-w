// Code generated by Stately. DO NOT EDIT.

package schema

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/StatelyCloud/go-sdk/pb/db"
	"github.com/StatelyCloud/go-sdk/stately"
)

// A "lease" gives users temporary access to a resource.
//
// Lease items can be accessed via the following key paths:
// * /user-:user_id/res-:resource_id/lease-:id
// * /res-:resource_id/lease-:id
// * /lease-:id
type Lease struct {
	// A unique identifier for the lease itself.
	Id uuid.UUID `protobuf:"bytes,1" json:"id,omitempty"`

	// The user that this lease is granted to.
	UserId uuid.UUID `protobuf:"bytes,2" json:"user_id,omitempty"`

	// The resource this lease grants access to.
	ResourceId uuid.UUID `protobuf:"bytes,3" json:"resource_id,omitempty"`

	// Allow the user to specify why they needed the lease.
	Reason string `protobuf:"bytes,4" json:"reason,omitempty"`

	// How long is this lease for? This is measured from when the lease was last modified.
	DurationSeconds time.Duration `protobuf:"zigzag64,5" json:"duration_seconds,omitempty,string"`

	// Last touch time allows us to extend a lease by updating it.
	LastTouched time.Time `protobuf:"zigzag64,6" json:"lastTouched,omitempty,string"`

	CreatedAt time.Time `protobuf:"zigzag64,7" json:"createdAt,omitempty,string"`

	// Who has approved this? The lease is not considered valid until approved by another person.
	Approver uuid.UUID `protobuf:"bytes,8" json:"approver,omitempty"`
}

// GetId is a nil-safe getter for field Id.
func (x *Lease) GetId() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.Id
}

// GetUserId is a nil-safe getter for field UserId.
func (x *Lease) GetUserId() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.UserId
}

// GetResourceId is a nil-safe getter for field ResourceId.
func (x *Lease) GetResourceId() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.ResourceId
}

// GetReason is a nil-safe getter for field Reason.
func (x *Lease) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}

// GetDurationSeconds is a nil-safe getter for field DurationSeconds.
func (x *Lease) GetDurationSeconds() time.Duration {
	if x == nil {
		return 0
	}
	return x.DurationSeconds
}

// GetLastTouched is a nil-safe getter for field LastTouched.
func (x *Lease) GetLastTouched() time.Time {
	if x == nil {
		return time.Time{}
	}
	return x.LastTouched
}

// GetCreatedAt is a nil-safe getter for field CreatedAt.
func (x *Lease) GetCreatedAt() time.Time {
	if x == nil {
		return time.Time{}
	}
	return x.CreatedAt
}

// GetApprover is a nil-safe getter for field Approver.
func (x *Lease) GetApprover() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.Approver
}

// MarshalJSON implements a custom JSON marshaller for Lease.
func (x Lease) MarshalJSON() ([]byte, error) {
	type Alias Lease
	aux := &struct {
		*Alias
		Id              []byte `json:"id,omitempty"`
		UserId          []byte `json:"user_id,omitempty"`
		ResourceId      []byte `json:"resource_id,omitempty"`
		DurationSeconds int64  `json:"duration_seconds,omitempty,string"`
		LastTouched     int64  `json:"lastTouched,omitempty,string"`
		CreatedAt       int64  `json:"createdAt,omitempty,string"`
		Approver        []byte `json:"approver,omitempty"`
	}{
		Alias:           (*Alias)(&x),
		Id:              uuidToBinary(x.Id),
		UserId:          uuidToBinary(x.UserId),
		ResourceId:      uuidToBinary(x.ResourceId),
		DurationSeconds: int64(x.DurationSeconds.Seconds()),
		LastTouched:     int64(x.LastTouched.UnixMilli()),
		CreatedAt:       int64(x.CreatedAt.UnixMilli()),
		Approver:        uuidToBinary(x.Approver),
	}
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler for Lease.
func (x *Lease) UnmarshalJSON(data []byte) error {
	type Alias Lease
	aux := &struct {
		*Alias
		Id              []byte `json:"id,omitempty"`
		UserId          []byte `json:"user_id,omitempty"`
		ResourceId      []byte `json:"resource_id,omitempty"`
		DurationSeconds int64  `json:"duration_seconds,omitempty,string"`
		LastTouched     int64  `json:"lastTouched,omitempty,string"`
		CreatedAt       int64  `json:"createdAt,omitempty,string"`
		Approver        []byte `json:"approver,omitempty"`
	}{Alias: (*Alias)(x)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	x.Id = binaryToUUID(aux.Id)
	x.UserId = binaryToUUID(aux.UserId)
	x.ResourceId = binaryToUUID(aux.ResourceId)
	x.DurationSeconds = time.Duration(aux.DurationSeconds) * time.Second
	x.LastTouched = time.UnixMilli(int64(aux.LastTouched))
	x.CreatedAt = time.UnixMilli(int64(aux.CreatedAt))
	x.Approver = binaryToUUID(aux.Approver)
	return nil
}

// StatelyItemType is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Lease) StatelyItemType() string {
	return "Lease"
}

// UnmarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Lease) UnmarshalStately(item *db.Item) error {
	return x.Unmarshal(item.GetProto())
}

// MarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Lease) MarshalStately() (*db.Item, error) {
	return marshalStatelyItem(x, x.StatelyItemType())
}

// KeyPath constructs and returns the primary key for this ItemType,
// based on the template `/user-:user_id/res-:resource_id/lease-:id` defined in schema.
// Note: The key constructed here will only be valid if the required key fields are set.
func (x *Lease) KeyPath() string {
	return "/user-" + stately.ToKeyID([16]byte(x.GetUserId())) +
		"/res-" + stately.ToKeyID([16]byte(x.GetResourceId())) +
		"/lease-" + stately.ToKeyID([16]byte(x.GetId()))
}

// A system is a resource that users can access.
//
// Resource items can be accessed via the following key paths:
// * /res-:id
type Resource struct {
	Id uuid.UUID `protobuf:"bytes,1" json:"id,omitempty"`

	Name string `protobuf:"bytes,2" json:"name,omitempty"`

	CreatedAt time.Time `protobuf:"zigzag64,3" json:"createdAt,omitempty,string"`
}

// GetId is a nil-safe getter for field Id.
func (x *Resource) GetId() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.Id
}

// GetName is a nil-safe getter for field Name.
func (x *Resource) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}

// GetCreatedAt is a nil-safe getter for field CreatedAt.
func (x *Resource) GetCreatedAt() time.Time {
	if x == nil {
		return time.Time{}
	}
	return x.CreatedAt
}

// MarshalJSON implements a custom JSON marshaller for Resource.
func (x Resource) MarshalJSON() ([]byte, error) {
	type Alias Resource
	aux := &struct {
		*Alias
		Id        []byte `json:"id,omitempty"`
		CreatedAt int64  `json:"createdAt,omitempty,string"`
	}{
		Alias:     (*Alias)(&x),
		Id:        uuidToBinary(x.Id),
		CreatedAt: int64(x.CreatedAt.UnixMilli()),
	}
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler for Resource.
func (x *Resource) UnmarshalJSON(data []byte) error {
	type Alias Resource
	aux := &struct {
		*Alias
		Id        []byte `json:"id,omitempty"`
		CreatedAt int64  `json:"createdAt,omitempty,string"`
	}{Alias: (*Alias)(x)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	x.Id = binaryToUUID(aux.Id)
	x.CreatedAt = time.UnixMilli(int64(aux.CreatedAt))
	return nil
}

// StatelyItemType is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Resource) StatelyItemType() string {
	return "Resource"
}

// UnmarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Resource) UnmarshalStately(item *db.Item) error {
	return x.Unmarshal(item.GetProto())
}

// MarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *Resource) MarshalStately() (*db.Item, error) {
	return marshalStatelyItem(x, x.StatelyItemType())
}

// KeyPath constructs and returns the primary key for this ItemType,
// based on the template `/res-:id` defined in schema.
// Note: The key constructed here will only be valid if the required key fields are set.
func (x *Resource) KeyPath() string {
	return "/res-" + stately.ToKeyID([16]byte(x.GetId()))
}

// A basic User object
//
// User items can be accessed via the following key paths:
// * /user-:id
// * /user_email-:email
type User struct {
	Id uuid.UUID `protobuf:"bytes,1" json:"id,omitempty"`

	DisplayName string `protobuf:"bytes,2" json:"displayName,omitempty"`

	Email string `protobuf:"bytes,3" json:"email,omitempty"`

	CreatedAt time.Time `protobuf:"zigzag64,4" json:"createdAt,omitempty,string"`
}

// GetId is a nil-safe getter for field Id.
func (x *User) GetId() uuid.UUID {
	if x == nil {
		return uuid.Nil
	}
	return x.Id
}

// GetDisplayName is a nil-safe getter for field DisplayName.
func (x *User) GetDisplayName() string {
	if x == nil {
		return ""
	}
	return x.DisplayName
}

// GetEmail is a nil-safe getter for field Email.
func (x *User) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}

// GetCreatedAt is a nil-safe getter for field CreatedAt.
func (x *User) GetCreatedAt() time.Time {
	if x == nil {
		return time.Time{}
	}
	return x.CreatedAt
}

// MarshalJSON implements a custom JSON marshaller for User.
func (x User) MarshalJSON() ([]byte, error) {
	type Alias User
	aux := &struct {
		*Alias
		Id        []byte `json:"id,omitempty"`
		CreatedAt int64  `json:"createdAt,omitempty,string"`
	}{
		Alias:     (*Alias)(&x),
		Id:        uuidToBinary(x.Id),
		CreatedAt: int64(x.CreatedAt.UnixMilli()),
	}
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler for User.
func (x *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
		Id        []byte `json:"id,omitempty"`
		CreatedAt int64  `json:"createdAt,omitempty,string"`
	}{Alias: (*Alias)(x)}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	x.Id = binaryToUUID(aux.Id)
	x.CreatedAt = time.UnixMilli(int64(aux.CreatedAt))
	return nil
}

// StatelyItemType is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *User) StatelyItemType() string {
	return "User"
}

// UnmarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *User) UnmarshalStately(item *db.Item) error {
	return x.Unmarshal(item.GetProto())
}

// MarshalStately is part of the stately.Item interface which is used by the golang SDK.
// For usage, please refer to the stately.Item interface documentation.
func (x *User) MarshalStately() (*db.Item, error) {
	return marshalStatelyItem(x, x.StatelyItemType())
}

// KeyPath constructs and returns the primary key for this ItemType,
// based on the template `/user-:id` defined in schema.
// Note: The key constructed here will only be valid if the required key fields are set.
func (x *User) KeyPath() string {
	return "/user-" + stately.ToKeyID([16]byte(x.GetId()))
}

type marshallerIFace interface {
	Marshal() ([]byte, error)
}

func marshalStatelyItem(msg marshallerIFace, itemType string) (*db.Item, error) {
	b, err := msg.Marshal()
	if err != nil {
		return nil, err
	}
	return &db.Item{
		Payload: &db.Item_Proto{
			Proto: b,
		},
		ItemType: itemType,
	}, nil
}

// TypeMapper defines a stately.ItemTypeMapper that unmarshals the wire format of your data
// into your SDK item types.
//
// Valid item types are:
// *Lease
// *Resource
// *User
func TypeMapper(item *db.Item) (stately.Item, error) {
	var result stately.Item
	switch item.ItemType {
	case "Lease":
		result = &Lease{}
	case "Resource":
		result = &Resource{}
	case "User":
		result = &User{}
	default:
		return nil, stately.UnknownItemTypeError{item.ItemType}
	}
	if err := result.UnmarshalStately(item); err != nil {
		return nil, err
	}
	return result, nil
}

func mapSlice[Tin any, Tout any](s []Tin, convert func(Tin) Tout) []Tout {
	out := make([]Tout, len(s))
	for i, v := range s {
		out[i] = convert(v)
	}
	return out
}

func mapSliceErr[Tin any, Tout any](s []Tin, convert func(Tin) (Tout, error)) ([]Tout, error) {
	out := make([]Tout, len(s))
	var err error
	for i, v := range s {
		out[i], err = convert(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func uuidToBinary(id uuid.UUID) []byte {
	return id[:]
}

func binaryToUUID(b []byte) uuid.UUID {
	if len(b) == 0 {
		return uuid.Nil
	}
	return uuid.UUID(b)
}
