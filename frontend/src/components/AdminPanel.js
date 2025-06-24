import React, { useEffect, useState } from 'react';
import UserManagement from './UserManagement';

const API_URL = 'http://localhost:3000/api/admin/questions';

function AdminPanel({ onLogout }) {
  const [questions, setQuestions] = useState([]);
  const [token, setToken] = useState('');
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({ id: '', exam_id: '', question_text: '', correct_answer: '', type: 'pilihan_ganda', options: ['', '', '', ''] });
  const [isEdit, setIsEdit] = useState(false);
  const [notif, setNotif] = useState('');
  const [notifType, setNotifType] = useState('');
  const [activeTab, setActiveTab] = useState('questions');

  useEffect(() => {
    // Ambil token dari localStorage (asumsi login admin sudah simpan token)
    const t = localStorage.getItem('token');
    setToken(t || '');
    fetchQuestions(t);
  }, []);

  const fetchQuestions = async (t) => {
    setLoading(true);
    const res = await fetch(API_URL, {
      headers: { Authorization: 'Bearer ' + t }
    });
    if (res.ok) {
      const data = await res.json();
      setQuestions(data);
    }
    setLoading(false);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Hapus soal ini?')) return;
    try {
      const res = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
        headers: { Authorization: 'Bearer ' + token }
      });
      if (res.ok) {
        setNotif('Soal berhasil dihapus!');
        setNotifType('success');
      } else {
        setNotif('Gagal menghapus soal!');
        setNotifType('error');
      }
    } catch {
      setNotif('Terjadi error jaringan!');
      setNotifType('error');
    }
    fetchQuestions(token);
  };

  const handleChange = e => {
    const { name, value } = e.target;
    if (name.startsWith('option_')) {
      const idx = parseInt(name.split('_')[1], 10);
      const newOptions = [...form.options];
      newOptions[idx] = value;
      setForm({ ...form, options: newOptions });
    } else {
      setForm({ ...form, [name]: value });
    }
  };

  const handleSubmit = async e => {
    e.preventDefault();
    if (!form.exam_id || !form.question_text || !form.correct_answer || !form.type) {
      setNotif('Semua field wajib diisi!');
      setNotifType('error');
      return;
    }
    if (form.type === 'pilihan_ganda' && form.options.filter(opt => opt.trim()).length < 2) {
      setNotif('Minimal 2 opsi untuk pilihan ganda!');
      setNotifType('error');
      return;
    }
    const payload = {
      ...form,
      exam_id: parseInt(form.exam_id, 10),
      options: form.type === 'pilihan_ganda' ? form.options.filter(opt => opt.trim()) : undefined
    };
    try {
      let res;
      if (isEdit) {
        res = await fetch(`${API_URL}/${form.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
          body: JSON.stringify(payload)
        });
      } else {
        res = await fetch(API_URL, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json', Authorization: 'Bearer ' + token },
          body: JSON.stringify(payload)
        });
      }
      if (res.ok) {
        setNotif(isEdit ? 'Soal berhasil diupdate!' : 'Soal berhasil ditambah!');
        setNotifType('success');
      } else {
        setNotif('Gagal menyimpan soal!');
        setNotifType('error');
      }
    } catch {
      setNotif('Terjadi error jaringan!');
      setNotifType('error');
    }
    setForm({ id: '', exam_id: '', question_text: '', correct_answer: '', type: 'pilihan_ganda', options: ['', '', '', ''] });
    setIsEdit(false);
    fetchQuestions(token);
  };

  const handleEdit = q => {
    setForm({
      ...q,
      options: Array.isArray(q.options) && q.options.length ? q.options : ['', '', '', '']
    });
    setIsEdit(true);
  };

  const handleCancel = () => {
    setForm({ id: '', exam_id: '', question_text: '', correct_answer: '', type: 'pilihan_ganda', options: ['', '', '', ''] });
    setIsEdit(false);
  };

  const handleExport = () => {
    setNotif('Export hasil ujian dimulai...');
    setNotifType('success');
    window.open('http://localhost:3000/api/admin/export?token=' + token, '_blank');
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Panel Admin</h2>
        {onLogout && <button onClick={onLogout}>Logout</button>}
      </div>
      
      {/* Tab Navigation */}
      <div style={{ marginBottom: 20 }}>
        <button 
          onClick={() => setActiveTab('questions')}
          style={{
            marginRight: 10,
            backgroundColor: activeTab === 'questions' ? '#007bff' : '#f8f9fa',
            color: activeTab === 'questions' ? 'white' : 'black',
            border: '1px solid #dee2e6',
            padding: '8px 16px',
            cursor: 'pointer'
          }}
        >
          Kelola Soal
        </button>
        <button 
          onClick={() => setActiveTab('users')}
          style={{
            backgroundColor: activeTab === 'users' ? '#007bff' : '#f8f9fa',
            color: activeTab === 'users' ? 'white' : 'black',
            border: '1px solid #dee2e6',
            padding: '8px 16px',
            cursor: 'pointer'
          }}
        >
          Kelola User
        </button>
      </div>

      {/* Tab Content */}
      {activeTab === 'questions' ? (
        <div>
          <h3>{isEdit ? 'Edit Soal' : 'Tambah Soal'}</h3>
          {notif && <div style={{color: notifType === 'success' ? 'green' : 'red', marginBottom: 10}}>{notif}</div>}
          <form onSubmit={handleSubmit} style={{ marginBottom: 20 }}>
            <input name="exam_id" placeholder="Exam ID" value={form.exam_id} onChange={handleChange} required />{' '}
            <input name="question_text" placeholder="Pertanyaan" value={form.question_text} onChange={handleChange} required />{' '}
            <select name="type" value={form.type} onChange={handleChange} required style={{marginRight: 10}}>
              <option value="pilihan_ganda">Pilihan Ganda</option>
              <option value="isian">Isian Singkat</option>
            </select>
            {form.type === 'pilihan_ganda' && (
              <div style={{margin:'10px 0'}}>
                {(Array.isArray(form.options) ? form.options : ['', '', '', '']).map((opt, idx) => (
                  <input key={idx} name={`option_${idx}`} placeholder={`Opsi ${String.fromCharCode(65+idx)}`} value={opt} onChange={handleChange} style={{marginRight:5}} />
                ))}
                <select name="correct_answer" value={form.correct_answer} onChange={handleChange} required>
                  <option value="">Pilih Jawaban Benar</option>
                  {(Array.isArray(form.options) ? form.options : ['', '', '', '']).map((opt, idx) => opt.trim() && <option key={idx} value={opt}>{String.fromCharCode(65+idx)}. {opt}</option>)}
                </select>
              </div>
            )}
            {form.type === 'isian' && (
              <input name="correct_answer" placeholder="Jawaban Benar" value={form.correct_answer} onChange={handleChange} required style={{marginRight:5}} />
            )}
            <button type="submit">{isEdit ? 'Update' : 'Tambah'}</button>
            {isEdit && <button type="button" onClick={handleCancel}>Batal</button>}
          </form>
          <button onClick={handleExport} style={{marginBottom: 20}}>Export Hasil Ujian (CSV)</button>
          {loading ? <p>Loading...</p> : (
            <table border="1" cellPadding="8">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Exam ID</th>
                  <th>Pertanyaan</th>
                  <th>Jawaban Benar</th>
                  <th>Aksi</th>
                </tr>
              </thead>
              <tbody>
                {questions.map(q => (
                  <tr key={q.id}>
                    <td>{q.id}</td>
                    <td>{q.exam_id}</td>
                    <td>{q.question_text}</td>
                    <td>{q.correct_answer}</td>
                    <td>
                      <button onClick={() => handleEdit(q)}>Edit</button>{' '}
                      <button onClick={() => handleDelete(q.id)}>Hapus</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      ) : (
        <UserManagement />
      )}
    </div>
  );
}

export default AdminPanel; 