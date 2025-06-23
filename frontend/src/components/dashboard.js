import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const Dashboard = ({ user, onLogout }) => {
    const [exams, setExams] = useState([]);
    const navigate = useNavigate();

    useEffect(() => {
        axios.get('http://localhost:3000/api/exams')
            .then((res) => setExams(res.data))
            .catch((err) => console.error(err));
    }, []);

    const handleStartExam = (exam) => {
        navigate(`/exam/${exam.id}`);
    };

    return (
        <div className="dashboard-container">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2>Selamat datang, {user.email}!</h2>
                <button onClick={onLogout} style={{ height: 'fit-content' }}>Logout</button>
            </div>
            <h1>DASHBOARD PESERTA</h1>
            {exams.map((exam) => (
                <div key={exam.id} className="exam-card">
                    <h2>{exam.title}</h2>
                    <p>Durasi: {exam.duration} menit</p>
                    <button onClick={() => handleStartExam(exam)}>Mulai Ujian</button>
                </div>
            ))}
        </div>
    );
};

export default Dashboard;