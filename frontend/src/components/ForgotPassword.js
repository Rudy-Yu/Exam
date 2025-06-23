import React from 'react';
import { Link } from 'react-router-dom';

const ForgotPassword = () => {
    return (
        <div className="forgot-password-container">
            <h1>Lupa Password</h1>
            <p>Fitur untuk reset password akan diimplementasikan di sini.</p>
            <Link to="/login">Kembali ke Login</Link>
        </div>
    );
};

export default ForgotPassword; 