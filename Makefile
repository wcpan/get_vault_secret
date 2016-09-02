DEPS=\
  github.com/Sirupsen/logrus \
  github.com/aws/aws-sdk-go/aws \
  github.com/hashicorp/vault \
  github.com/mitchellh/gox


bootstrap:
	@for dep in $(DEPS) ; do \
		echo "Installing $$dep" ; \
		go get $$dep; \
	done



.PHONY: bootstrap

