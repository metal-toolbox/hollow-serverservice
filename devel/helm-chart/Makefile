.DEFAULT_GOAL := help

## install helm chart for the server service local development environment
local-devel: kubectl-ctx-kind
	helm install serverservice-devel . -f values.yaml

## upgrade helm chart for local devel environment
local-devel-upgrade: kubectl-ctx-kind
	helm upgrade serverservice-devel . -f values.yaml

## port forward hollow server service port (runs in foreground)
port-forward-hss: kubectl-ctx-kind
	kubectl port-forward deployment/serverservice 8000:8000

## port forward crdb service port (runs in foreground)
port-forward-crdb: kubectl-ctx-kind
	kubectl port-forward deployment/crdb 26257:26257

## connect to crdb with psql (requires port-forward-crdb)
psql-crdb: kubectl-ctx-kind
	psql -d "postgresql://root@localhost:26257/defaultdb?sslmode=disable"

## set kube ctx to kind cluster
kubectl-ctx-kind:
	export KUBECONFIG=~/.kube/config_kind
	kubectl config use-context kind-kind

# https://gist.github.com/prwhite/8168133
# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


TARGET_MAX_CHAR_NUM=20
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
