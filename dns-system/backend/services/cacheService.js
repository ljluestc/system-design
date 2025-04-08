const redis = require("redis");
const client = redis.createClient({ host: "redis", port: 6379 });

client.on("error", (err) => console.error("Redis error:", err));

async function getFromCache(key) {
  return new Promise((resolve, reject) => {
    client.get(key, (err, data) => {
      if (err) return reject(err);
      resolve(data ? JSON.parse(data) : null);
    });
  });
}

async function setInCache(key, value, ttl = 3600) {
  return new Promise((resolve, reject) => {
    client.setex(key, ttl, JSON.stringify(value), (err) => {
      if (err) return reject(err);
      resolve();
    });
  });
}

module.exports = { getFromCache, setInCache };