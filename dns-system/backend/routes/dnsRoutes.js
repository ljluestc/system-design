const express = require("express");
const { resolveDomain } = require("../controllers/dnsController");

const router = express.Router();

router.get("/resolve", resolveDomain);

module.exports = router;