USER_ROOT:=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

$(USER_ROOT)-codegen: tools-codegen
	$(BIN)/codegen -p cerulean.ir/modules/user/aaa
	$(BIN)/codegen -p cerulean.ir/modules/user/controllers/user

$(USER_ROOT)-migration: tools-go-bindata
	cd $(USER_ROOT)/migrations && $(BIN)/go-bindata -nometadata -o migration.gen.go -nomemcopy=true -pkg=migrations ./db/...

#$(USER_ROOT)-lint: $(LINTER)
#	$(LINTERCMD) $(USER_ROOT)/...

.PHONY: $(USER_ROOT)-codegen #$(USER_ROOT)-migration $(USER_ROOT)-lint