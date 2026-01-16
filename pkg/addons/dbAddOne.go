package addons

type DbAddOneConfig struct {
	ServiceName string
	Image       string
	Environment string
	Port        string
	Volume      string
	VolumeName  string
	DBName      string
	DBEnvPrefix string
	Import      string
	Driver      string
	DSN         string
	OutputFile  string
}

var DbRegistory = map[string]DbAddOneConfig{
	"postgres": {
		ServiceName: "postgres_bp",
		Image:       "postgres:latest",
		Environment: `      POSTGRES_DB: ${GONE_DB_DATABASE}
      POSTGRES_USER: ${GONE_DB_USERNAME}
      POSTGRES_PASSWORD: ${GONE_DB_PASSWORD}`,
		Port:        "5432",
		Volume:      "postgres_volume_bp:/var/lib/postgresql/data",
		VolumeName:  "postgres_volume_bp",
		DBName:      "PostgreSQL",
		DBEnvPrefix: "BLUEPRINT",
		Import:      `_ "github.com/jackc/pgx/v5/stdlib"`,
		Driver:      "pgx",
		DSN:         "postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
	},

	"mysql": {
		ServiceName: "mysql_bp",
		Image:       "mysql:8",
		Environment: `
      MYSQL_DATABASE: ${GONE_DB_DATABASE}
      MYSQL_USER: ${GONE_DB_USERNAME}
      MYSQL_PASSWORD: ${GONE_DB_PASSWORD}
      MYSQL_ROOT_PASSWORD: ${GONE_DB_PASSWORD}`,
		Port:        "3306",
		Volume:      "mysql_volume_bp:/var/lib/mysql",
		VolumeName:  "mysql_volume_bp",
		DBName:      "MySQL",
		DBEnvPrefix: "BLUEPRINT",
		Import:      `_ "github.com/go-sql-driver/mysql"`,
		Driver:      "mysql",
		DSN:         "%s:%s@tcp(%s:%s)/%s",
		OutputFile:  "mysql.go",
	},

	"mongodb": {
		ServiceName: "mongo_bp",
		Image:       "mongo:latest",
		Environment: `
      MONGO_INITDB_DATABASE: ${GONE_DB_DATABASE}
      MONGO_INITDB_ROOT_USERNAME: ${GONE_DB_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: ${GONE_DB_PASSWORD}`,
		Port:       "27017",
		Volume:     "mongo_volume_bp:/data/db",
		VolumeName: "mongo_volume_bp",
	},

	"sqlite": {
		ServiceName: "sqlite_bp",
		Image:       "alpine:latest",
		Environment: `
      # SQLite has no environment variables`,
		Port:        "0", // no port needed
		Volume:      "sqlite_volume_bp:/data",
		VolumeName:  "sqlite_volume_bp",
		DBName:      "SQLite",
		DBEnvPrefix: "BLUEPRINT",
		Import:      `_ "modernc.org/sqlite"`,
		Driver:      "sqlite",
		DSN:         "file:%s.db?_pragma=journal_mode(WAL)",
	},

	"cockroachdb": {
		ServiceName: "cockroach_bp",
		Image:       "cockroachdb/cockroach:latest",
		Environment: `
      # CockroachDB has no env vars needed in insecure mode`,
		Port:        "26257",
		Volume:      "cockroach_volume_bp:/cockroach/cockroach-data",
		VolumeName:  "cockroach_volume_bp",
		DBName:      "CockroachDB",
		DBEnvPrefix: "BLUEPRINT",
		Import:      `_ "github.com/jackc/pgx/v5/stdlib"`,
		Driver:      "pgx",
		DSN:         "postgres://%s:%s@%s:%s/%s?sslmode=disable",
	},

	"mariadb": {
		ServiceName: "mariadb_bp",
		Image:       "mariadb:latest",
		Environment: `
      MARIADB_DATABASE: ${GONE_DB_DATABASE}
      MARIADB_USER: ${GONE_DB_USERNAME}
      MARIADB_PASSWORD: ${GONE_DB_PASSWORD}
      MARIADB_ROOT_PASSWORD: ${GONE_DB_PASSWORD}`,
		Port:        "3306",
		Volume:      "mariadb_volume_bp:/var/lib/mysql",
		VolumeName:  "mariadb_volume_bp",
		DBName:      "MariaDB",
		DBEnvPrefix: "BLUEPRINT",
		Import:      `_ "github.com/go-sql-driver/mysql"`,
		Driver:      "mysql",
		DSN:         "%s:%s@tcp(%s:%s)/%s",
	},
}
