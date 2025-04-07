import React, { useState } from 'react';
import { login } from '../services/api';

function Login({ setAuth }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async () => {
    try {
      const { token, userId } = await login({ username, password });
      localStorage.setItem('token', token);
      localStorage.setItem('userId', userId);
      setAuth(true);
    } catch (error) {
      alert('Login failed');
    }
  };

  return (
    <div>
      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
      />
      <button onClick={handleLogin}>Login</button>
    </div>
  );
}

export default Login;