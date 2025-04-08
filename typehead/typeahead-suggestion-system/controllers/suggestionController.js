const trieService = require("../services/trieService");
const cacheService = require("../services/cacheService");
const logger = require("../utils/logger");

// Constants for configuration
const CACHE_TTL = 3600; // 1 hour in seconds
const MAX_SUGGESTIONS = 10;

/**
 * Handles GET requests to retrieve suggestions for a given query.
 * Integrates caching and trie-based suggestion generation.
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @returns {void}
 */
async function getSuggestions(req, res) {
  const { query } = req.query;

  try {
    // Log request start
    logger.log(`Processing suggestion request for query: ${query}`);

    // Input validation
    if (!query || typeof query !== "string") {
      logger.log("Invalid query parameter received");
      return res.status(400).json({
        error: "Query parameter is required and must be a string",
        timestamp: new Date().toISOString(),
      });
    }

    if (query.length < 1) {
      logger.log("Query too short, returning empty suggestions");
      return res.status(400).json({
        suggestions: [],
        message: "Query must be at least 1 character",
      });
    }

    // Sanitize query
    const sanitizedQuery = query.trim().toLowerCase();
    logger.log(`Sanitized query: ${sanitizedQuery}`);

    // Check cache first
    let suggestions;
    try {
      suggestions = await cacheService.get(sanitizedQuery);
      if (suggestions) {
        logger.log(`Cache hit for query: ${sanitizedQuery}`);
        return res.json({
          suggestions,
          source: "cache",
          timestamp: new Date().toISOString(),
        });
      }
    } catch (cacheError) {
      logger.log(`Cache retrieval error: ${cacheError.message}`);
    }

    // Fetch from trie if cache miss
    logger.log(`Cache miss, fetching from trie for: ${sanitizedQuery}`);
    try {
      suggestions = trieService.getSuggestions(sanitizedQuery);
      logger.log(`Trie returned ${suggestions.length} suggestions`);
    } catch (trieError) {
      logger.log(`Trie error: ${trieError.message}`);
      return res.status(500).json({
        error: "Failed to fetch suggestions from trie",
        details: trieError.message,
      });
    }

    // Limit suggestions
    const limitedSuggestions = suggestions.slice(0, MAX_SUGGESTIONS);
    logger.log(`Limited suggestions to ${limitedSuggestions.length}`);

    // Cache the results
    try {
      await cacheService.set(sanitizedQuery, limitedSuggestions, CACHE_TTL);
      logger.log(`Cached suggestions for: ${sanitizedQuery}`);
    } catch (cacheSetError) {
      logger.log(`Cache set error: ${cacheSetError.message}`);
    }

    // Send response
    res.json({
      suggestions: limitedSuggestions,
      source: "trie",
      timestamp: new Date().toISOString(),
    });
    logger.log("Suggestion request completed successfully");
  } catch (error) {
    logger.log(`Unexpected error: ${error.message}`);
    res.status(500).json({
      error: "Internal server error",
      details: error.message,
      timestamp: new Date().toISOString(),
    });
  }
}

module.exports = { getSuggestions };

/*
  Expansion Notes:
  - Add extensive JSDoc (10-15 lines per function with examples).
  - Implement detailed logging for every step (e.g., entry, validation, cache hit/miss).
  - Add multiple validation checks (e.g., length, special characters).
  - Include retry logic for cache and trie operations.
  - Add performance metrics (e.g., response time calculation).
  - Implement fallback mechanisms (e.g., static suggestions if trie fails).
  - Add comments explaining each block (5-10 lines per section).
  - Repeat error handling blocks with slight variations.
  - Add mock personalization logic (e.g., user-specific filtering).
*/