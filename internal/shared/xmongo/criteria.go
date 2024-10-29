package xmongo

import "go.mongodb.org/mongo-driver/bson"

// Criteria represents a criteria to be used in a query to get events from the storage.
type Criteria interface {
	// ToBSON returns the criteria as a BSON document.
	ToBSON() bson.D
}

// CriteriaBuilder is a function that builds a Criteria.
type CriteriaBuilder func() Criteria

// idCriteria is a Criteria to get events by ID.
type idCriteria struct {
	id string
}

// WithIDCriteria returns a CriteriaBuilder that builds a Criteria to get events by ID.
func WithEventIDCriteria(id string) CriteriaBuilder {
	return func() Criteria {
		return &idCriteria{id: id}
	}
}

// ToBSON returns the criteria as a BSON document.
func (c *idCriteria) ToBSON() bson.D {
	return bson.D{{Key: "_id", Value: c.id}}
}

// aggregateIDCriteria is a Criteria to get events by aggregate ID.
type aggregateIDCriteria struct {
	aggregateID string
}

// WithAggregateIDCriteria returns a CriteriaBuilder that builds a Criteria to get events by aggregate ID.
func WithAggregateIDCriteria(aggregateID string) CriteriaBuilder {
	return func() Criteria {
		return &aggregateIDCriteria{aggregateID: aggregateID}
	}
}

// ToBSON returns the criteria as a BSON document.
func (c *aggregateIDCriteria) ToBSON() bson.D {
	return bson.D{{Key: "aggregate_id", Value: c.aggregateID}}
}

// aggregateTypeCriteria is a Criteria to get events by aggregate type.
type aggregateTypeCriteria struct {
	aggregateType string
}

// WithAggregateTypeCriteria returns a CriteriaBuilder that builds a Criteria to get events by aggregate type.
func WithAggregateTypeCriteria(aggregateType string) CriteriaBuilder {
	return func() Criteria {
		return &aggregateTypeCriteria{aggregateType: aggregateType}
	}
}

// ToBSON returns the criteria as a BSON document.
func (c *aggregateTypeCriteria) ToBSON() bson.D {
	return bson.D{{Key: "aggregate_type", Value: c.aggregateType}}
}

// aggregateVersionCriteria is a Criteria to get events by aggregate version.
type aggregateVersionCriteria struct {
	aggregateVersion int
}

// WithAggregateVersionCriteria returns a CriteriaBuilder that builds a Criteria to get events by aggregate version.
func WithAggregateVersionCriteria(aggregateVersion int) CriteriaBuilder {
	return func() Criteria {
		return &aggregateVersionCriteria{aggregateVersion: aggregateVersion}
	}
}

// ToBSON returns the criteria as a BSON document.
func (c *aggregateVersionCriteria) ToBSON() bson.D {
	return bson.D{{Key: "aggregate_version", Value: c.aggregateVersion}}
}

// typeCriteria is a Criteria to get events by type.
type typeCriteria struct {
	eventType string
}

// WithEventTypeCriteria returns a CriteriaBuilder that builds a Criteria to get events by type.
func WithEventTypeCriteria(eventType string) CriteriaBuilder {
	return func() Criteria {
		return &typeCriteria{eventType: eventType}
	}
}

// ToBSON returns the criteria as a BSON document.
func (c *typeCriteria) ToBSON() bson.D {
	return bson.D{{Key: "type", Value: c.eventType}}
}

// TODO: Implement versionCriteria
// versionCriteria is a Criteria to get events by version.
// type versionCriteria struct {
// 	version int
// }

// func WithEventVersionCriteria(version int) CriteriaBuilder {
// 	return func() Criteria {
// 		return &versionCriteria{version: version}
// 	}
// }

// // ToBSON returns the criteria as a BSON document.
// func (c *versionCriteria) ToBSON() bson.D {
// 	return bson.D{{Key: "version", Value: c.version}}
// }

// And is a CriteriaBuilder that builds a Criteria that is the result of the logical AND operation between two Criteria.
func And(crs ...Criteria) CriteriaBuilder {
	return func() Criteria {
		return &andCriteria{criterias: crs}
	}
}

// andCriteria is a Criteria that is the result of the logical AND operation between two Criteria.
type andCriteria struct {
	criterias []Criteria
}

// ToBSON returns the criteria as a BSON document.
func (c *andCriteria) ToBSON() bson.D {
	bsonArray := make(bson.A, len(c.criterias))
	for i, c := range c.criterias {
		bsonArray[i] = c.ToBSON()
	}
	return bson.D{{Key: "$and", Value: bsonArray}}
}

// Or is a CriteriaBuilder that builds a Criteria that is the result of the logical OR operation between two Criteria.
func Or(crs ...Criteria) CriteriaBuilder {
	return func() Criteria {
		return &orCriteria{criterias: crs}
	}
}

// orCriteria is a Criteria that is the result of the logical OR operation between two Criteria.
type orCriteria struct {
	criterias []Criteria
}

// ToBSON returns the criteria as a BSON document.
func (c *orCriteria) ToBSON() bson.D {
	bsonArray := make(bson.A, len(c.criterias))
	for i, c := range c.criterias {
		bsonArray[i] = c.ToBSON()
	}
	return bson.D{{Key: "$or", Value: bsonArray}}
}

// Not is a CriteriaBuilder that builds a Criteria that is the result of the logical NOT operation of a Criteria.
func Not(crs ...Criteria) CriteriaBuilder {
	return func() Criteria {
		return &notCriteria{criterias: crs}
	}
}

// notCriteria is a Criteria that is the result of the logical NOT operation of a Criteria.
type notCriteria struct {
	criterias []Criteria
}

// ToBSON returns the criteria as a BSON document.
func (c *notCriteria) ToBSON() bson.D {
	bsonArray := make(bson.A, len(c.criterias))
	for i, c := range c.criterias {
		bsonArray[i] = c.ToBSON()
	}
	return bson.D{{Key: "$nor", Value: bsonArray}}
}
