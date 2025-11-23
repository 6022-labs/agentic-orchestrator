export enum DatabaseConnectionDriver {
  None = "none", // Used for the default value
  MySQL = "mysql",
  MariaDB = "mariadb",
  Postgres = "postgres",
}

export const DatabaseConnectionDriverDisplayNames: Record<
  DatabaseConnectionDriver,
  string
> = {
  [DatabaseConnectionDriver.None]: "",
  [DatabaseConnectionDriver.MySQL]: "MySQL",
  [DatabaseConnectionDriver.MariaDB]: "MariaDB",
  [DatabaseConnectionDriver.Postgres]: "Postgres",
};

export function getAllDatabaseConnectionDrivers(): DatabaseConnectionDriver[] {
  return [
    DatabaseConnectionDriver.MySQL,
    DatabaseConnectionDriver.MariaDB,
    DatabaseConnectionDriver.Postgres,
  ];
}
