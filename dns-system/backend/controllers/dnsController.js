const dns = require("dns").promises;
const { getFromCache, setInCache } = require("../services/cacheService");
const logger = require("../utils/logger");

async function resolveDomain(req, res) {
  const { domain } = req.query;

  if (!domain) {
    return res.status(400).json({ error: "Domain is required" });
  }

  try {
    let ip = await getFromCache(domain);
    if (ip) {
      logger.info(`Cache hit for ${domain}: ${ip}`);
      return res.json({ ip });
    }

    ip = (await dns.resolve4(domain))[0];
    await setInCache(domain, ip);
    logger.info(`Resolved and cached ${domain}: ${ip}`);
    res.json({ ip });
  } catch (error) {
    logger.error(`DNS resolution failed for ${domain}: ${error.message}`);
    res.status(500).json({ error: "DNS resolution failed" });
  }
}

module.exports = { resolveDomain };