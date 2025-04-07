async function generateResponse(message) {
    // Simulate AI response generation (replace with real API/model in production)
    await new Promise((resolve) => setTimeout(resolve, 1000));
    return `Response to: ${message}`;
  }
  
  module.exports = { generateResponse };