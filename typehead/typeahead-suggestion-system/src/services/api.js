import axios from "axios";

// Configuration constants
const API_BASE_URL = "http://localhost:3001/api";
const TIMEOUT = 5000;
const MAX_RETRIES = 3;

/**
 * API client instance with base configuration.
 */
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: TIMEOUT,
});

/**
 * Fetches suggestions from the backend API with retry logic.
 * @param {string} query - The search query
 * @returns {Promise<Object>} The API response with suggestions
 * @throws {Error} If all retries fail
 */
export async function getSuggestions(query) {
  let retries = 0;

  while (retries < MAX_RETRIES) {
    try {
      console.log(`[API] Fetching suggestions for: ${query}, attempt ${retries + 1}`);
      const res = await api.get(`/suggestions?query=${query}`);
      console.log(`[API] Successfully fetched: ${res.data.suggestions.length} suggestions`);
      return res.data;
    } catch (error) {
      retries++;
      console.error(`[API] Error on attempt ${retries}: ${error.message}`);
      if (retries === MAX_RETRIES) {
        throw new Error(`Failed to fetch suggestions after ${MAX_RETRIES} attempts: ${error.message}`);
      }
      await new Promise((resolve) => setTimeout(resolve, 1000 * retries)); // Exponential backoff
    }
  }
}