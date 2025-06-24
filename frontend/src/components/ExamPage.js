import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';

const API_URL = process.env.NODE_ENV === 'production'
    ? (process.env.REACT_APP_API_URL || 'https://api.yourdomain.com')
    : 'http://localhost:3000';

const ExamPage = () => {
    const { id } = useParams();
    const [questions, setQuestions] = useState([]);
    const [answers, setAnswers] = useState({});
    const [score, setScore] = useState(null);
    const [timeLeft, setTimeLeft] = useState(null);
    const [notif, setNotif] = useState('');
    const [notifType, setNotifType] = useState('');
    const [examStarted, setExamStarted] = useState(false);

    // Start exam and get server time
    useEffect(() => {
        const startExam = async () => {
            try {
                const res = await axios.post(`${API_URL}/api/exam/${id}/start`, {}, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                if (res.data.success) {
                    setExamStarted(true);
                    setTimeLeft(res.data.duration);
                }
            } catch (err) {
                setNotif('Gagal memulai ujian');
                setNotifType('error');
            }
        };

        startExam();
    }, [id]);

    // Fetch questions
    useEffect(() => {
        if (!examStarted) return;

        axios.get(`${API_URL}/api/exam/${id}/questions`, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
        })
            .then(res => setQuestions(res.data))
            .catch(err => {
                setNotif('Gagal mengambil soal');
                setNotifType('error');
            });
    }, [id, examStarted]);

    // Timer sync with server
    useEffect(() => {
        if (!examStarted) return;

        const timer = setInterval(async () => {
            try {
                const res = await axios.get(`${API_URL}/api/exam/${id}/timer`, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                });
                if (res.data.success) {
                    setTimeLeft(res.data.remaining_time);
                    if (res.data.remaining_time <= 0) {
                        handleSubmit();
                    }
                }
            } catch (err) {
                console.error('Error syncing timer:', err);
            }
        }, 1000);

        return () => clearInterval(timer);
    }, [examStarted, id]);

    // Auto-save answers every 30 seconds
    useEffect(() => {
        if (!examStarted || Object.keys(answers).length === 0) return;

        const autoSave = setInterval(() => {
            Object.entries(answers).forEach(([questionId, answerText]) => {
                axios.post(`${API_URL}/api/answers/draft`, {
                    question_id: questionId,
                    answer_text: answerText
                }, {
                    headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
                })
                    .then(() => {
                        setNotif('Jawaban tersimpan otomatis');
                        setNotifType('success');
                        setTimeout(() => setNotif(''), 2000);
                    })
                    .catch(err => console.error('Auto-save error:', err));
            });
        }, 30000);

        return () => clearInterval(autoSave);
    }, [answers, examStarted]);

    const handleAnswerChange = (questionId, answer) => {
        setAnswers(prev => ({ ...prev, [questionId]: answer }));
    };

    const handleSubmit = async () => {
        if (!window.confirm('Yakin ingin mengumpulkan jawaban?')) return;

        try {
            // Submit semua jawaban sebagai final
            const answersArray = Object.entries(answers).map(([questionId, answerText]) => ({
                question_id: parseInt(questionId),
                answer_text: answerText
            }));

            await axios.post(`${API_URL}/api/answers/submit`, answersArray, {
                headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
            });

            setNotif('Jawaban berhasil dikumpulkan');
            setNotifType('success');
            
            // Hitung skor
            let currentScore = 0;
            questions.forEach(q => {
                if (answers[q.id] === q.correct_answer) {
                    currentScore += q.weight || 1; // Gunakan bobot soal jika ada
                }
            });
            setScore(currentScore);
        } catch (err) {
            setNotif('Gagal mengumpulkan jawaban');
            setNotifType('error');
        }
    };

    if (score !== null) {
        return (
            <div className="exam-page">
                <h1>Hasil Ujian</h1>
                <h2>Skor Anda: {score} dari {questions.reduce((acc, q) => acc + (q.weight || 1), 0)}</h2>
                <Link to="/dashboard">Kembali ke Dashboard</Link>
            </div>
        );
    }

    const formatTime = (seconds) => {
        if (seconds === null) return '--:--';
        const minutes = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${minutes}:${secs < 10 ? '0' : ''}${secs}`;
    };

    return (
        <div className="exam-page">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h1>Ujian #{id}</h1>
                <h2>Sisa Waktu: {formatTime(timeLeft)}</h2>
            </div>
            
            {notif && (
                <div style={{
                    color: notifType === 'success' ? 'green' : 'red',
                    marginBottom: 10,
                    padding: '10px',
                    borderRadius: '4px',
                    backgroundColor: notifType === 'success' ? '#e8f5e9' : '#ffebee'
                }}>
                    {notif}
                </div>
            )}

            {questions.map((q, index) => (
                <div key={q.id} style={{ marginBottom: '20px' }}>
                    <p>
                        <b>{index + 1}. {q.text}</b>
                        {q.weight > 1 && <span style={{ color: '#666', fontSize: '0.9em' }}> (Bobot: {q.weight})</span>}
                    </p>
                    {q.options.map(option => (
                        <div key={option}>
                            <input
                                type="radio"
                                name={`question-${q.id}`}
                                value={option}
                                checked={answers[q.id] === option}
                                onChange={() => handleAnswerChange(q.id, option)}
                            />
                            <label>{option}</label>
                        </div>
                    ))}
                </div>
            ))}

            <div style={{ marginTop: '20px' }}>
                <button 
                    onClick={handleSubmit}
                    style={{
                        backgroundColor: '#4CAF50',
                        color: 'white',
                        padding: '10px 20px',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer'
                    }}
                >
                    Submit Jawaban
                </button>
                <Link 
                    to="/dashboard" 
                    style={{ 
                        marginLeft: '10px',
                        color: '#666',
                        textDecoration: 'none'
                    }}
                >
                    Kembali ke Dashboard
                </Link>
            </div>
        </div>
    );
};

export default ExamPage; 