import React, { useState, useEffect } from 'react';
import axios from 'axios';

const Exam = () => {
    const [questions, setQuestions] = useState([]);
    const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);

    useEffect(() => {
        axios.get('http://localhost:3000/api/questions')
            .then((res) => setQuestions(res.data))
            .catch((err) => console.error(err));
    }, []);

    const handleNextQuestion = () => {
        setCurrentQuestionIndex(currentQuestionIndex + 1);
    };

    return (
        <div className="exam-container">
            <h1>SOAL UJIAN</h1>
            <div className="question">
                <h2>Soal {currentQuestionIndex + 1}: {questions[currentQuestionIndex]?.question_text}</h2>
                <ul>
                    {questions[currentQuestionIndex]?.options?.map((option, index) => (
                        <li key={index}>{option}</li>
                    ))}
                </ul>
            </div>
            <button onClick={handleNextQuestion}>Selanjutnya</button>
        </div>
    );
};

export default Exam;