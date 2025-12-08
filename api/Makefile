.PHONY: generate
generate:
	bunx tsp compile .
	oapi-codegen -package=api -generate=chi-server,models,strict-server -o gen.go tsp-output/schema/openapi.yaml