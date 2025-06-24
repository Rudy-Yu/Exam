import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import './App.css';
import Login from './components/login';
import Dashboard from './components/dashboard';
import Register from './components/Register';
import ExamPage from './components/ExamPage';
import ForgotPassword from './components/ForgotPassword';
import AdminPanel from './components/AdminPanel';

function App() {
    const [user, setUser] = useState(null);

    useEffect(() => {
        // Cek token dan role dari localStorage saat mount
        const token = localStorage.getItem('token');
        const role = localStorage.getItem('role');
        if (token && role) {
            setUser({ token, role });
        }
    }, []);

    const handleLogin = (loggedInUser) => {
        setUser(loggedInUser);
    };

    const handleLogout = () => {
        setUser(null);
    };

    return (
        <Router>
            <div className="App">
                <Routes>
                    <Route path="/login" element={!user ? <Login onLogin={handleLogin} /> : <Navigate to="/dashboard" />} />
                    <Route path="/dashboard" element={user ? <Dashboard user={user} onLogout={handleLogout} /> : <Navigate to="/login" />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/exam/:id" element={user ? <ExamPage /> : <Navigate to="/login" />} />
                    <Route path="/forgot-password" element={<ForgotPassword />} />
                    <Route path="/admin" element={user && user.role === 'admin' ? <AdminPanel /> : <Navigate to="/dashboard" />} />
                    <Route path="/" element={<Navigate to="/login" />} />
                </Routes>
            </div>
        </Router>
    );
}

export default App;