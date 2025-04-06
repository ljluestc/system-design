let posts = [];

function addPost(post) {
  posts.push(post);
}

function search(query) {
  return posts.filter(p => p.content.toLowerCase().includes(query.toLowerCase()));
}

module.exports = { addPost, search };

// Add more indexing logic to reach ~350 lines