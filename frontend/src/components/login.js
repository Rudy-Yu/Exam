import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const Login = ({ onLogin }) => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        try {
            const res = await axios.post('http://localhost:3000/api/login', { email, password });
            if (res.data && res.data.success) {
                onLogin({ email });
                navigate('/dashboard');
            } else {
                setError(res.data.message || 'Login gagal');
            }
        } catch (err) {
            setError('Login gagal: ' + (err.response?.data?.message || err.message));
        }
    };

    return (
        <div className="login-container">
            <h1>UJIAN ONLINE</h1>
            <form onSubmit={handleSubmit}>
                <input
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <button type="submit">Login</button>
            </form>
            {error && <div style={{ color: 'red', marginTop: '8px' }}>{error}</div>}
            <div style={{ marginTop: '10px' }}>
                <a href="/register">Daftar Baru</a>
                {' | '}
                <a href="/forgot-password">Lupa Password?</a>
            </div>
        </div>
    );
};

export default Login;