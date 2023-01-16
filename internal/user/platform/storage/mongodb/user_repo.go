package mongodb

type UserRepository struct {
}

func NewUserRepositoryMongo() *UserRepository {
	return &UserRepository{}
}
