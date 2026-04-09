data "external_schema" "gorm-mysql" {
  program = [
    "go",
    "run",
    "./loader",
    "-dialect",
    "mysql",
  ]
}

data "external_schema" "gorm-postgres" {
  program = [
    "go",
    "run",
    "./loader",
    "-dialect",
    "postgres",
  ]
}

data "external_schema" "gorm-mariadb" {
  program = [
    "go",
    "run",
    "./loader",
    "-dialect",
    "mysql",
  ]
}

env "gorm-mysql" {
  src = data.external_schema.gorm-mysql.url
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://mysql"
    format = goose
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "gorm-postgres" {
  src = data.external_schema.gorm-postgres.url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://postgres"
    format = goose
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "gorm-mariadb" {
  src = data.external_schema.gorm-mariadb.url
  dev = "docker://maria/10.7/test"
  migration {
    dir = "file://mariadb"
    format = goose
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}