CREATE TABLE exams (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    exam_id INT REFERENCES exams(id),
    question_text TEXT NOT NULL,
    correct_answer VARCHAR(255) NOT NULL
);

CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    participant_id INT,
    question_id INT REFERENCES questions(id),
    answer_text VARCHAR(255) NOT NULL,
    submitted_at TIMESTAMP DEFAULT NOW()
);