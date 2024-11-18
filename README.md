# Rock Paper Sockets
This is a simple web application to play Rock Paper Scissors with other people, using web sockets.

## Building the Project
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

In another terminal instance, you can run the following to setup a simple HTTP server for the HTML code:
```sh
python -m http.server 5000
```

You can use any port number you want for the HTTP server except for 8000, as `uvicorn` will run on 8000.