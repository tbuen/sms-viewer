BINARY=smsv
ifneq ($(RELEASE),)
	export CGO_CPPFLAGS=$(CPPFLAGS)
	export CGO_CFLAGS=$(CFLAGS)
	export CGO_CXXFLAGS=$(CXXFLAGS)
	export CGO_LDFLAGS=$(LDFLAGS)
	FLAGS=-v -tags release -buildmode=pie -trimpath -mod=readonly -modcacherw -ldflags='-linkmode=external -w -X=github.com/tbuen/sms-viewer/internal/app.version=$(RELEASE) -extldflags=$(LDFLAGS)'
else
	PREBUILD=fmt
	FLAGS=-v
endif

.PHONY: $(BINARY)
$(BINARY): $(PREBUILD)
	go build $(FLAGS) ./cmd/...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	@rm -rf smsv
	@rm -rf build/package/sms-viewer* build/package/pkg build/package/src

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/bin
	cp $(BINARY) $(DESTDIR)/usr/bin/
