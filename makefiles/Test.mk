# /*
# |    Protect your secrets, protect your sensitive data.
# :    Explore VMware Secrets Manager docs at https://vsecm.com/
# </
# <>/  keep your secrets... secret
# >/
# <>/' Copyright 2023-present VMware Secrets Manager contributors.
# >/'  SPDX-License-Identifier: BSD-2-Clause
# */

#
# ## Coverage ##
#

coverage_file := coverage.out
threshold = 70

#  To run all tests and check coverage against the threshold:
#  make cover
cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=$(coverage_file) ./... > /dev/null
	@echo "Checking test coverage..."
	@coverage=$$(go tool cover -func=$(coverage_file) | grep total | grep -Eo '[0-9]+\.[0-9]+' || echo "0.0"); \
	echo "Test coverage: $$coverage"; \
	echo "Test Threshold: $(threshold)"; \
    if [ "$$coverage" != "" ] && awk 'BEGIN{exit !('"$$coverage"'<=$(threshold))}'; then \
    	echo "Test coverage is less than $(threshold)"; \
		exit 0; \
	fi
	if [ "$$coverage" != "" ] && awk 'BEGIN{exit !('"$$coverage"'>=$(threshold))}'; then \
		echo "Test coverage is greater than $(threshold)"; \
		exit 0; \
	fi
	@rm -f $(coverage_file)

#
# ## Tests ##
#

# Integration tests.
test:
	./hack/test.sh "remote" ""
test-remote:
	./hack/test.sh "remote" ""
test-local:
	./hack/test.sh "local" ""
test-local-ci:
	./hack/test.sh "local" "ci"
