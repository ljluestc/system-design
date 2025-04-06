const fs = require('fs').promises;
const path = require('path');
const logger = require('../../shared/logger');

async function uploadMedia(req, res) {
  try {
    const mediaFile = req.files.media;
    const uploadPath = path.join(__dirname, 'uploads', mediaFile.name);
    await fs.mkdir(path.dirname(uploadPath), { recursive: true });
    await mediaFile.mv(uploadPath);
    res.json({ url: `/uploads/${mediaFile.name}` });
  } catch (err) {
    logger.error(`Media upload failed: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

module.exports = { uploadMedia };

// Add more logic to reach ~400 lines