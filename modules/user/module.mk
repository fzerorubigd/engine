USER_ROOT:=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

$(USER_ROOT)-swagger: swagger-to-go
	$(BIN)/swagger-to-go -pkg userpb -file $(USER_ROOT)/proto/user.swagger.json > $(USER_ROOT)/proto/user.swagger.pb.go


$(USER_ROOT)-migration: $(BIN)/go-bindata
	cd $(USER_ROOT)/migrations && $(BIN)/go-bindata -nometadata -o migration.gen.go -nomemcopy=true -pkg=migrations ./db/...

#$(USER_ROOT)-lint: $(LINTER)
#	$(LINTERCMD) $(USER_ROOT)/...

.PHONY: $(USER_ROOT)-codegen #$(USER_ROOT)-migration $(USER_ROOT)-lint