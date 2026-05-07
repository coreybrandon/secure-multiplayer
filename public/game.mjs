import Player from './Player.mjs';
import Collectible from './Collectible.mjs';

const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');

const MOVE_DIST = 5;
const PLAYER_SIZE = 30;
const COLLECTIBLE_RADIUS = 10;

let myId = null;
let players = {};
let collectible = null;
const keysDown = new Set();

const socket = new WebSocket(`ws://${location.host}/ws`);

socket.addEventListener('message', (event) => {
  const msg = JSON.parse(event.data);

  if (msg.type === 'init') {
    myId = msg.player.id;
    players = {};
    msg.players.forEach(p => { players[p.id] = new Player(p); });
    collectible = msg.collectible ? new Collectible(msg.collectible) : null;
  } else if (msg.type === 'state') {
    players = {};
    msg.players.forEach(p => { players[p.id] = new Player(p); });
    collectible = msg.collectible ? new Collectible(msg.collectible) : null;
  } else if (msg.type === 'player-left') {
    delete players[msg.id];
  }
});

socket.addEventListener('close', () => {
  console.log('disconnected from server');
});

document.addEventListener('keydown', (e) => keysDown.add(e.key));
document.addEventListener('keyup', (e) => keysDown.delete(e.key));

function sendMove(dir) {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ type: 'move', dir, dist: MOVE_DIST }));
  }
}

function processInput() {
  if (!myId) return;
  if (keysDown.has('ArrowUp')    || keysDown.has('w')) sendMove('up');
  if (keysDown.has('ArrowDown')  || keysDown.has('s')) sendMove('down');
  if (keysDown.has('ArrowLeft')  || keysDown.has('a')) sendMove('left');
  if (keysDown.has('ArrowRight') || keysDown.has('d')) sendMove('right');
}

function draw() {
  ctx.fillStyle = '#1a1a2e';
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  if (collectible) {
    ctx.fillStyle = '#ffd700';
    ctx.beginPath();
    ctx.arc(
      collectible.x + COLLECTIBLE_RADIUS,
      collectible.y + COLLECTIBLE_RADIUS,
      COLLECTIBLE_RADIUS, 0, Math.PI * 2
    );
    ctx.fill();
    ctx.fillStyle = '#ffffff';
    ctx.font = '7px "Press Start 2P"';
    ctx.fillText(`+${collectible.value}`, collectible.x, collectible.y - 2);
  }

  Object.values(players).forEach(p => {
    const isMe = p.id === myId;
    ctx.fillStyle = isMe ? '#e94560' : '#0f3460';
    ctx.fillRect(p.x, p.y, PLAYER_SIZE, PLAYER_SIZE);
    ctx.fillStyle = '#ffffff';
    ctx.font = '7px "Press Start 2P"';
    ctx.fillText(String(p.score), p.x + 2, p.y - 4);
  });

  if (myId && players[myId]) {
    const rank = players[myId].calculateRank(Object.values(players));
    ctx.fillStyle = '#ffffff';
    ctx.font = '10px "Press Start 2P"';
    ctx.fillText(rank, 10, 20);
  }
}

let lastMoveTime = 0;
const MOVE_INTERVAL = 1000 / 30;

function gameLoop(timestamp) {
  if (timestamp - lastMoveTime >= MOVE_INTERVAL) {
    lastMoveTime = timestamp;
    processInput();
  }
  draw();
  requestAnimationFrame(gameLoop);
}

requestAnimationFrame(gameLoop);
