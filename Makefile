test:
	go test -race ./...

serve:
	docker run --rm -it --name echo-docs -v ${CURDIR}/website:/home/app -w /home/app -p 3000:3000 -u node node:lts /bin/bash -c "npm install && npm start -- --host=0.0.0.0"

.PHONY: test serve
