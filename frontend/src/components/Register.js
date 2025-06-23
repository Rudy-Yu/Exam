import React, { useState } from 'react';
import { Link } from 'react-router-dom';

const Register = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        // Logika pendaftaran akan ditambahkan di sini
        setMessage('Pendaftaran berhasil! Silakan login.');
    };

    return (
        <div className="register-container">
            <h1>DAFTAR BARU</h1>
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
                <button type="submit">Daftar</button>
            </form>
            {message && <div style={{ color: 'green', marginTop: '10px' }}>{message}</div>}
            <div style={{ marginTop: '10px' }}>
                <Link to="/login">Sudah punya akun? Login</Link>
            </div>
        </div>
    );
};

export default Register; 