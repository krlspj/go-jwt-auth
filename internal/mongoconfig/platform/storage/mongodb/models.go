package mongodb

import (
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Refresh   string             `bson:"refresh,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
	//CreatedAt primitive.Timestamp `bson:"created_at,omitempty"`
}

func (m mongoConfig) toDomain() (domain.Config, error) {
	config := new(domain.Config)
	config.SetID(m.ID.Hex())
	//b, err := strconv.ParseBool(m.Refresh)
	//if err != nil {
	//	return domain.Config{}, err
	//}
	//config.SetRefresh(b)
	config.SetRefresh(m.Refresh)
	config.SetCreatedAt(m.CreatedAt)

	return *config, nil
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
	//mConfig.Refresh = strconv.FormatBool(dConfig.Refresh())
	mConfig.Refresh = dConfig.Refresh()
	mConfig.CreatedAt = dConfig.CreatedAt()
	//mConfig.CreatedAt = primitive.Timestamp{T: dConfig.CreatedAt()}

	return *mConfig, nil
}
