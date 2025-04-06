const axios = require('axios');

async function get(url) {
  return axios.get(url);
}

async function post(url, data) {
  return axios.post(url, data);
}

module.exports = { get, post };

// Add more HTTP utilities to reach ~150 lines