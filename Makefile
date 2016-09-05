DEPS=\
  github.com/Sirupsen/logrus \
  github.com/aws/aws-sdk-go/aws \
  github.com/hashicorp/vault \
  github.com/mitchellh/gox \
  github.com/tcnksm/ghr

bootstrap:
	@for dep in $(DEPS) ; do \
		echo "Installing $$dep" ; \
		go get $$dep; \
	done

build:
	./build.sh
	
release:
	ghr -replace v0.0.1 pkg

.PHONY: bootstrap
