async function generateResponse(message) {
    // Mock implementation; integrate with Hugging Face or OpenAI API in production
    await new Promise((resolve) => setTimeout(resolve, 500));
    return `AI says: ${message}`;
  }
  
  module.exports = { generateResponse };