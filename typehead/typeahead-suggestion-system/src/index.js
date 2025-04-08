import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import "./styles/App.css";

/**
 * Entry point for the Typeahead Suggestion System React application.
 * Renders the App component into the DOM and handles initial setup.
 */
function initializeApp() {
  try {
    // Render the application in strict mode for development checks
    ReactDOM.render(
      <React.StrictMode>
        <App />
      </React.StrictMode>,
      document.getElementById("root")
    );
    console.log("React application rendered successfully");
  } catch (error) {
    console.error("Failed to render React app:", error);
    document.body.innerHTML = "<h1>Error loading application</h1>";
  }
}

// Execute initialization
initializeApp();

// Add basic environment check
if (process.env.NODE_ENV === "development") {
  console.log("Running in development mode");
}