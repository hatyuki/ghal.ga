.PHONY: deploy

deploy:
	umask 022
	apex deploy -s GITHUB_TOKEN=$(GITHUB_TOKEN)
