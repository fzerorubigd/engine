MISC_ROOT:=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

$(MISC_ROOT)-swagger: swagger-to-go
	$(BIN)/swagger-to-go -pkg miscpb -file $(MISC_ROOT)/proto/misc.swagger.json > $(MISC_ROOT)/proto/misc.swagger.pb.go


#$(MISC_ROOT)-migration: $(BIN)/go-bindata
#	cd $(MISC_ROOT)/migrations && $(BIN)/go-bindata -nometadata -o migration.gen.go -nomemcopy=true -pkg=migrations ./db/...

#$(USER_ROOT)-lint: $(LINTER)
#	$(LINTERCMD) $(MISC_ROOT)/...

.PHONY: $(MISC_ROOT)-codegen #$(MISC_ROOT)-migration $(MISC_ROOT)-lint