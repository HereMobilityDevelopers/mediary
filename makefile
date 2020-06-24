.travis:                # Travis CI (see .travis.yml), runs tests
ifndef TRAVIS
	@echo "Fail: requires Travis runtime"
else
	@$(MAKE) test --no-print-directory && goveralls -coverprofile=./coverage.out -service=travis-ci
endif
