
.PHONY: proto
proto:
	@echo "--> Generating proto bindings..."
	@buf --config tools/buf/buf.yaml --template tools/buf/buf.gen.yaml generate