const mongoose = require("mongoose");

const ResourceRecordSchema = new mongoose.Schema({
  type: { type: String, required: true, enum: ["A", "NS", "CNAME", "MX"] },
  name: { type: String, required: true },
  value: { type: String, required: true },
  ttl: { type: Number, default: 3600 },
  createdAt: { type: Date, default: Date.now },
});

module.exports = mongoose.model("ResourceRecord", ResourceRecordSchema);