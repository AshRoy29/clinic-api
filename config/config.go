package config

type Config struct {
	Server   Server
	Database Database
	Jwt      Jwt
}

type Server struct {
	Port string
}

type Database struct {
	Uri          string
	DatabaseName string
	Username     string
	Password     string
}

type Jwt struct {
	Secret string
}

func (c *Config) Read() {

	c.Server.Port = ":4000"
	c.Database.Uri = "mongodb://localhost:27017"
	c.Database.DatabaseName = "clinic"
	c.Database.Username = ""
	c.Database.Password = ""
	//c.Jwt.Secret = "secretkeyjwt"
	//flag.StringVar(&c.Server.Port, "Server", ":4000", "Server Port")
	////flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	//flag.StringVar(&c.Database.Uri, "Uri", "mongodb://localhost:27017", "MongoDB connection string")
	//flag.StringVar(&c.Database.DatabaseName, "DBname", "clinic", "MongoDB database name string")
	//flag.StringVar(&c.Database.Username, "Username", "", "MongoDB username string")
	//flag.StringVar(&c.Database.Password, "Password", "", "MongoDB password string")
	//flag.StringVar(&c.Jwt.Secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "secret")

}
