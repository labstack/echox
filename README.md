# Echo Extra

### Running Recipes

- Go into any cookbook folder and run `go run server.go`

### Running/Developing Website Locally

- Ensure that you have Node.js installed on your computer. If you don't have it, download and install it from the
  official Node.js website (https://nodejs.org). Or you could use NVM (https://github.com/nvm-sh/nvm) to use appropriate
  version of Node.
- Inside the `website` folder, install the required dependencies by running the following command `npm install`.
- Run the website locally by running the following command `npm start`.

```bash
cd website
npm install
npm start
```

### Running/Developing Website in docker

This will serve website on http://localhost:3000/

```bash
docker run --rm -it --name echo-docs -v ${PWD}/website:/home/app -w /home/app -p 3000:3000 -u node node:lts /bin/bash -c "npm install && npm start -- --host=0.0.0.0"
```
