# Rock Paper Sockets
This is a simple web application to play Rock Paper Scissors with other people, using web sockets.

## Building the Project - Python
This project uses [FastAPI](https://fastapi.tiangolo.com/). To install the various packages, run the following in the terminal:
```sh
pip install fastapi uvicorn
```

Additionally, make sure `websockets` is up to date:
```sh
pip install --upgrade websockets
```

Then to run the application, execute the following:
```sh
uvicorn main:app --reload
```

Then you can visit `localhost:8000` in your browser to view the webpage.

## Building the Project - Golang
This project uses [Fiber](https://docs.gofiber.io/). Ensure you have at least Golang version 1.18.
To install the various packages, run the following in the terminal:
```sh
go mod tidy
```

Then to run the application, execute the following:
```sh
go run .
```

You can also optionally install `air`, which allows for live reloading the server by running the following:
```sh
go install github.com/cosmtrek/air@latest
```

Then to run the application using `air`, simply type:
```sh
air
```