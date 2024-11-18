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

Then you can visit `localhost:8000` in your browser to view the webpage.