import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';

const ExamPage = () => {
    const { id } = useParams();
    const [questions, setQuestions] = useState([]);
    const [answers, setAnswers] = useState({});
    const [score, setScore] = useState(null);
    const [timeLeft, setTimeLeft] = useState(600); // 10 menit dalam detik

    useEffect(() => {
        axios.get(`http://localhost:3000/api/exam/${id}/questions`)
            .then(res => setQuestions(res.data))
            .catch(err => console.error(err));
    }, [id]);

    useEffect(() => {
        if (score === null) {
            const timer = setInterval(() => {
                setTimeLeft(prevTime => {
                    if (prevTime <= 1) {
                        clearInterval(timer);
                        handleSubmit();
                        return 0;
                    }
                    return prevTime - 1;
                });
            }, 1000);

            return () => clearInterval(timer);
        }
    }, [score]);

    const handleAnswerChange = (questionId, answer) => {
        setAnswers(prev => ({ ...prev, [questionId]: answer }));
    };

    const handleSubmit = () => {
        let currentScore = 0;
        questions.forEach(q => {
            if (answers[q.id] === q.correct_answer) {
                currentScore++;
            }
        });
        setScore(currentScore);
    };

    if (score !== null) {
        return (
            <div className="exam-page">
                <h1>Hasil Ujian</h1>
                <h2>Skor Anda: {score} dari {questions.length}</h2>
                <Link to="/dashboard">Kembali ke Dashboard</Link>
            </div>
        );
    }

    const formatTime = (seconds) => {
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
            {questions.map((q, index) => (
                <div key={q.id} style={{ marginBottom: '20px' }}>
                    <p><b>{index + 1}. {q.text}</b></p>
                    {q.options.map(option => (
                        <div key={option}>
                            <input
                                type="radio"
                                name={`question-${q.id}`}
                                value={option}
                                onChange={() => handleAnswerChange(q.id, option)}
                            />
                            <label>{option}</label>
                        </div>
                    ))}
                </div>
            ))}
            <div style={{ marginTop: '20px' }}>
                <button onClick={handleSubmit}>Submit Jawaban</button>
                <Link to="/dashboard" style={{ marginLeft: '10px' }}>Kembali ke Dashboard</Link>
            </div>
        </div>
    );
};

export default ExamPage; 