import React from "react";
import DNSResolver from "./components/DNSResolver";
import ErrorBoundary from "./components/ErrorBoundary";

function App() {
  return (
    <div className="app">
      <h1>DNS Resolver</h1>
      <ErrorBoundary>
        <DNSResolver />
      </ErrorBoundary>
    </div>
  );
}

export default App;