data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./loader",
  ]
}
env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://."
    format = goose
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
