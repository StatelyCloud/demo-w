// Code generated by Stately. DO NOT EDIT.

package schema

import (
	"context"
	"github.com/StatelyCloud/go-sdk/stately"
)

// NewClient is a convenient wrapper around stately.NewClient which creates a new client for the schema package
// while ensuring it uses the correct stately.ItemTypeMapper
func NewClient(ctx context.Context, storeID uint64, options ...*stately.Options) (stately.Client, error) {
	return stately.NewClient(ctx, storeID, 3, 4291558376530788, TypeMapper, options...)
}
