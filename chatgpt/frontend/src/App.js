import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import Login from "./components/Login";
import Register from "./components/Register";
import ChatWindow from "./components/ChatWindow";
import { authenticate } from "./services/api";
import "./styles/App.css";

/**
 * Main application component that handles routing and authentication state.
 * This component serves as the entry point for the React application,
 * managing user authentication and directing users to appropriate views.
 */

/* Constants */
const APP_NAME = "ChatGPT Clone";
const LOADING_TIMEOUT = 5000; // 5 seconds timeout for loading

/* Utility Functions */
const log = (message) => {
  console.log(`[${APP_NAME}] ${message}`);
};

/* Main App Component */
function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Effect to check authentication status on mount
  useEffect(() => {
    log("Initializing application...");
    const token = localStorage.getItem("token");
    if (token) {
      log("Found token, verifying...");
      authenticate(token)
        .then((valid) => {
          log("Token valid, setting authenticated state");
          setIsAuthenticated(valid);
          setLoading(false);
        })
        .catch((err) => {
          log(`Authentication failed: ${err.message}`);
          setIsAuthenticated(false);
          setLoading(false);
          setError("Failed to authenticate. Please log in again.");
        });
    } else {
      log("No token found, user not authenticated");
      setLoading(false);
    }

    // Cleanup function
    return () => log("Cleaning up App component");
  }, []);

  // Handle loading state with timeout
  useEffect(() => {
    const timer = setTimeout(() => {
      if (loading) {
        setLoading(false);
        setError("Loading timeout exceeded");
        log("Loading timeout exceeded");
      }
    }, LOADING_TIMEOUT);
    return () => clearTimeout(timer);
  }, [loading]);

  // Render loading state
  if (loading) {
    return (
      <div className="loading-container">
        <h2>Loading {APP_NAME}...</h2>
        <p>Please wait while we set things up.</p>
      </div>
    );
  }

  // Render error state
  if (error) {
    return (
      <div className="error-container">
        <h2>Error</h2>
        <p>{error}</p>
        <button onClick={() => window.location.reload()}>Retry</button>
      </div>
    );
  }

  // Main render with routing
  return (
    <Router>
      <div className="app">
        <header>
          <h1>{APP_NAME}</h1>
          <p>A real-time chat application powered by AI</p>
        </header>
        <main>
          <Switch>
            <Route path="/login">
              <Login setAuth={setIsAuthenticated} />
            </Route>
            <Route path="/register">
              <Register />
            </Route>
            <Route path="/">
              {isAuthenticated ? (
                <ChatWindow />
              ) : (
                <Login setAuth={setIsAuthenticated} />
              )}
            </Route>
          </Switch>
        </main>
        <footer>
          <p>&copy; 2023 {APP_NAME}. All rights reserved.</p>
        </footer>
      </div>
    </Router>
  );
}

/* Additional Components */
function WelcomeMessage() {
  return (
    <div className="welcome">
      <p>Welcome to {APP_NAME}! Start chatting now.</p>
    </div>
  );
}

/* Error Boundary Component */
class ErrorBoundary extends React.Component {
  state = { hasError: false };

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  render() {
    if (this.state.hasError) {
      return <h2>Something went wrong.</h2>;
    }
    return this.props.children;
  }
}

/* Export with Error Boundary */
export default function WrappedApp() {
  return (
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  );
}

// Add more comments and documentation
/**
 * Notes:
 * - This component uses React Router for navigation.
 * - Authentication is managed via a token stored in localStorage.
 * - Expand this file further by adding:
 *   - Context API for global state management.
 *   - More utility functions for analytics or logging.
 *   - Additional routes for user profiles or settings.
 */