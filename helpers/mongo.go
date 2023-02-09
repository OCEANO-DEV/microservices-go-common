package helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID primitive.ObjectID

// String convert ID to string.
func (id ID) String() string {
	return primitive.ObjectID(id).Hex()
}

// StringToID converts a string to ID.
func StringToID(s string) ID {
	_id, _ := primitive.ObjectIDFromHex(s)
	return ID(_id)
}

// IsValidID checks if ID is valid.
func IsValidID(s string) bool {
	return primitive.IsValidObjectID(s)
}

// NewID create a new id
func NewID() ID {
	return ID(primitive.NewObjectID())
}
