test:
	go test -race ./...

serve:
	cd website && hugo serve

.PHONY: test serve
