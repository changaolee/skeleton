# ==============================================================================
# 用来进行代码生成的 Makefile

.PHONY: gen.run
gen.run: gen.clean gen.docgo.doc

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name '*_generated.go' -delete

.PHONY: gen.docgo.doc
gen.docgo.doc:
	@echo "===========> Generating missing doc.go for go packages"
	@$(ROOT_DIR)/scripts/gendoc.sh

.PHONY: gen.docgo.check
gen.docgo.check: gen.docgo.doc
	@n="$$(git ls-files --others '*/doc.go' | wc -l)"; \
	if test "$$n" -gt 0; then \
		git ls-files --others '*/doc.go' | sed -e 's/^/  /'; \
		echo "$@: untracked doc.go file(s) exist in working directory" >&2 ; \
		false ; \
	fi

.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))

.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@$(ROOT_DIR)/scripts/gencerts.sh generate-skt-cert $(OUTPUT_DIR)/cert $(CA)
