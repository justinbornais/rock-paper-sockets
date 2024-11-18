let ws = null;
let gameCode = '';
let playerId = null;
let move = null;
let scores = null;

// Helper to display messages
function updateStatus(message) {
    document.getElementById("status").innerHTML = message;
}

function connectToGame() {
  gameCode = document.getElementById("game-code").value;
  if (!gameCode || gameCode.length !== 4) {
    updateStatus("Please enter a valid 4-digit game code.");
    return;
  }

  updateStatus('Waiting for another player...');

  ws = new WebSocket(`ws://localhost:8000/ws/${gameCode}`);
  ws.onopen = () => {
    updateStatus("Connected to game room.");
    document.getElementById("game-setup").style.display = "none";
    document.getElementById("game-play").style.display = "block";
    document.getElementById("room-code").innerText = `Game Code: ${gameCode}`;
  };

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);

    if (data.action === "game_started") {
      updateStatus("Game has started! Please make your move.");
      playerId = parseInt(data.player_id);
      document.getElementById("move-buttons").style.display = "block";
    } else if (data.action === "result") {
      const winner = data.winner;
      const opponentMove = data.opponent_move;
      scores = data.scores;

      document.getElementById("round-result").innerText =
        winner === null
          ? `It's a tie! Opponent played ${opponentMove}.`
          : `You ${winner === playerId ? "win" : "lose"} this round! You played ${move}. Opponent played ${opponentMove}.`;

      if (scores) {
        document.getElementById("round-result").innerText += ` Scores: You (${scores[playerId]}), Opponent (${scores[1 - playerId]})`;
      }
    } else if (data.action === "game_over") {
      const gameWinner = data.winner;
      document.getElementById("game-result").innerText =
        gameWinner === playerId
          ? "You win the game!"
          : "You lose the game!";
      document.getElementById("move-buttons").style.display = "none";
    }
  };

  ws.onclose = () => {
    updateStatus("Disconnected from game room.");
    document.getElementById("game-setup").style.display = "block";
    document.getElementById("game-play").style.display = "none";
  };

  ws.onerror = () => {
    updateStatus("Error connecting to the game room.");
  };
}

function sendMove(move) {
  if (!ws || ws.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({ action: "move", move }));
  document.getElementById("round-result").innerText = "Waiting for opponent's move...";
}

// Event listeners
document.getElementById("connect-btn").addEventListener("click", connectToGame);

document.querySelectorAll(".move-btn").forEach((button) => {
  button.addEventListener("click", () => {
    move = button.getAttribute("data-move");
    sendMove(move);
  });
});