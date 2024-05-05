default: build_container

build_container:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-common-hdwallet-api)
	$(eval context=$(or $(context),k0s-dev-cluster))
	$(eval platform=$(or $(platform),linux/amd64))

	$(eval short_commit_id=$(shell git rev-parse --short HEAD))
	$(eval commit_id=$(shell git rev-parse HEAD))
	$(eval build_number=0)
	$(eval build_date=$(shell date +%s))
	$(eval release_tag=$(shell git describe --tags $(commit_id))-$(short_commit_id)-$(build_number))

	docker build \
		--ssh default=$(SSH_AUTH_SOCK) \
		--platform $(platform) \
		--build-arg RELEASE_TAG=$(release_tag) \
		--build-arg COMMIT_ID=$(commit_id) \
		--build-arg SHORT_COMMIT_ID=$(short_commit_id) \
		--build-arg BUILD_NUMBER=$(build_number) \
		--build-arg BUILD_DATE_TS=$(build_date) \
		--tag $(container_registry):latest \
		--tag $(container_registry):$(build_tag) .

	docker push $(container_registry):$(build_tag)

.PHONY: build_container
