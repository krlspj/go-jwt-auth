package mongodb

import (
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoConfig struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Refresh   bool                `bson:"refresh,omitempty"`
	CreatedAt primitive.Timestamp `bson:"created_at,omitempty"`
}

func (m mongoConfig) toDomain() domain.Config {
	config := new(domain.Config)
	config.SetID(m.ID.Hex())
	//config.SetRefresh(strconv.FormatBool(m.Refresh))
	config.SetRefresh(m.Refresh)
	config.SetCreatedAt(m.CreatedAt.T)

	return *config
}

func toMongoConfig(dConfig domain.Config) (mongoConfig, error) {
	var oid primitive.ObjectID
	if dConfig.ID() != "" {
		var err error
		oid, err = primitive.ObjectIDFromHex(dConfig.ID())
		if err != nil {
			return mongoConfig{}, err
		}
	}
	mConfig := new(mongoConfig)
	mConfig.ID = oid
	//if dConfig.Refresh() != "" {
	//	boolValue, err := strconv.ParseBool(dConfig.Refresh())
	//	if err != nil {
	//		return mongoConfig{}, err
	//	}
	//	mConfig.Refresh = boolValue
	//}
	mConfig.Refresh = dConfig.Refresh()
	mConfig.CreatedAt = primitive.Timestamp{T: dConfig.CreatedAt()}

	return *mConfig, nil
}
