class Player {
  constructor({ id, x, y, score }) {
    this.id = id;
    this.x = x;
    this.y = y;
    this.score = score;
  }

  movePlayer(dir, dist) {
    if (dir === 'up')    this.y -= dist;
    if (dir === 'down')  this.y += dist;
    if (dir === 'left')  this.x -= dist;
    if (dir === 'right') this.x += dist;
  }

  collision(item) {
    return Math.abs(this.x - item.x) < 30 && Math.abs(this.y - item.y) < 30;
  }

  calculateRank(arr) {
    const sorted = [...arr].sort((a, b) => b.score - a.score);
    const rank = sorted.findIndex(p => p.id === this.id) + 1;
    return `Rank: ${rank} / ${arr.length}`;
  }
}

export default Player;
