.PHONY: check.lint

setup-git-hooks:
	$(info Setting up git hooks)
	@printf '#!/bin/sh \nmake pre-push-hook' > .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push

pre-push-hook: check.terraform check.lint

check.terraform:
	terraform fmt --check

check.lint:
	cd resize-go; make lint
