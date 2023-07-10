package datasource

import (
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func init() {
	print("asd")
}

func TestCreateIndexesWithValidIndexesShouldNotPanic(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("test name", func(mt *mtest.T) {
		var mongoDatasource = MongoDatasourceImpl{
			client: mt.Client,
		}
		mongoDatasource.CreateMongoIndexes()
	})
}
