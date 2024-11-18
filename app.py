from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from typing import Dict, List

# Set up variables.
app = FastAPI()
games = {}
WINNING_POINTS = 3

class GameRoom:
  def __init__(self, code):
    self.code = code
    self.players: List[WebSocket] = []
    self.moves = {}  # Maps each player to their move.
    self.scores = {0: 0, 1: 0}

  def add_player(self, websocket: WebSocket):
    if len(self.players) < 2:
      self.players.append(websocket)

  def all_moves_received(self):
    return len(self.moves) == 2

  def calculate_winner(self):
    moves = list(self.moves.values())
    if moves[0] == moves[1]:
      return None  # Tie.
    if (moves[0], moves[1]) in [("r", "s"), ("s", "p"), ("p", "r")]:
      return 0  # Player 1 wins.
    return 1  # Player 2 wins.

  def reset_moves(self):
    self.moves = {}

# Endpoint for the web socket.
@app.websocket("/ws/{game_code}")
async def websocket_endpoint(websocket: WebSocket, game_code: str):
  await websocket.accept()
  
  if game_code not in games:
    games[game_code] = GameRoom(game_code)

  game = games[game_code]

  # Check game availability.
  if len(game.players) >= 2:
    await websocket.send_text("Game room is full.")
    await websocket.close()
    return

  game.add_player(websocket)
  player_id = len(game.players) - 1

  try:
    while True:
      data = await websocket.receive_json()
      action = data.get("action")

      if action == "move":
        move = data.get("move")
        game.moves[player_id] = move

        # Calculate winner if both players made their move.
        if game.all_moves_received():
          winner = game.calculate_winner()
          if winner is not None:
            game.scores[winner] += 1

          for i, player in enumerate(game.players):
            await player.send_json({
              "action": "result",
              "winner": winner,
              "scores": game.scores,
              "opponent_move": game.moves[1 - i],
            })

          # Check if game is over.
          if max(game.scores.values()) == WINNING_POINTS:
            for player in game.players:
              await player.send_json({
                  "action": "game_over",
                  "winner": winner,
              })
            del games[game_code]
            break

          game.reset_moves()

  # Remove game in the event of a disconnect.
  except WebSocketDisconnect:
    game.players.remove(websocket)
    if not game.players:
      del games[game_code]