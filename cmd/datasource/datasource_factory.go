package datasource

import (
	"oauth/internal/config"
)

type Driver string

const (
	Mongo      Driver = "Mongo"
	SqlServer         = "SqlServer"
	Mysql             = "Mysql"
	Postgresql        = "Postgresql"
	Test              = "Test"
)

var mongoDatasource MongoDatasourceImpl

var Get = func() Datasource {
	var driver = Driver(config.GetProperty("Mongo"))
	switch driver {
	case Mongo:
		return mongoDatasource
	case Mysql:
		return nil
	case SqlServer:
		return nil
	case Postgresql:
		return nil
	case Test:
		return nil
	}
	return nil
}

var GetMongoDatasource = func() MongoDatasource {
	return mongoDatasource
}

var GetByDriver = func(driver Driver) any {
	switch driver {
	case Mongo:
		return mongoDatasource
	case Mysql:
		return nil
	case SqlServer:
		return nil
	case Postgresql:
		return nil
	case Test:
		return nil
	}
	return nil
}

var Setup = func() {
	var driver = Driver(config.GetProperty("DRIVER"))
	switch driver {
	case Mongo:
		mongoDatasource.SetupMongo()
	case Mysql:
	case SqlServer:
	case Postgresql:
	case Test:
	}
}

var SetupByDriver = func(driver Driver) {
	switch driver {
	case Mongo:
		mongoDatasource.SetupMongo()
	case Mysql:
	case SqlServer:
	case Postgresql:
	case Test:
	}
}
