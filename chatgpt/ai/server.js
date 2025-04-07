const express = require('express');
const { generateResponse } = require('./generate');

const app = express();
app.use(express.json());

app.post('/generate', async (req, res) => {
  const { message } = req.body;
  const response = await generateResponse(message);
  res.json({ response });
});

app.listen(3002, () => console.log('AI Service running on port 3002'));