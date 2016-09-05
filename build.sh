#!/bin/sh
XC_OSARCH=${XC_OSARCH:-"linux/amd64 darwin/amd64"}
gox -osarch="${XC_OSARCH}" -tags="get_vault_secret" -output "pkg/get-vault-secret_{{.OS}}_{{.Arch}}"
