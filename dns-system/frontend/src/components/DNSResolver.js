import React, { useState } from "react";
import { resolveDomain } from "../services/api";

function DNSResolver() {
  const [domain, setDomain] = useState("");
  const [ip, setIp] = useState("");
  const [error, setError] = useState("");

  const handleResolve = async () => {
    try {
      setError("");
      setIp("");
      const response = await resolveDomain(domain);
      setIp(response.ip);
    } catch (err) {
      setError("Failed to resolve domain. Please try again.");
    }
  };

  return (
    <div>
      <input
        type="text"
        value={domain}
        onChange={(e) => setDomain(e.target.value)}
        placeholder="Enter domain (e.g., google.com)"
      />
      <button onClick={handleResolve}>Resolve</button>
      {ip && <p>Resolved IP: {ip}</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
}

export default DNSResolver;