import React, { useEffect, useState } from 'react';

const API_URL = 'http://localhost:3000/api/admin/users';

function UserManagement() {
  const [users, setUsers] = useState([]);
  const [token, setToken] = useState('');
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({ id: '', email: '', role: 'user' });
  const [isEdit, setIsEdit] = useState(false);
  const [notif, setNotif] = useState('');
  const [notifType, setNotifType] = useState('');

  useEffect(() => {
    const t = localStorage.getItem('token');
    setToken(t || '');
    fetchUsers(t);
  }, []);

  const fetchUsers = async (t) => {
    setLoading(true);
    try {
      const res = await fetch(API_URL, {
        headers: { Authorization: 'Bearer ' + t }
      });
      if (res.ok) {
        const data = await res.json();
        setUsers(data);
      } else {
        setNotif('Gagal mengambil data user!');
        setNotifType('error');
      }
    } catch (error) {
      setNotif('Terjadi error jaringan!');
      setNotifType('error');
    }
    setLoading(false);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Hapus user ini?')) return;
    try {
      const res = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
        headers: { Authorization: 'Bearer ' + token }
      });
      if (res.ok) {
        setNotif('User berhasil dihapus!');
        setNotifType('success');
        fetchUsers(token);
      } else {
        setNotif('Gagal menghapus user!');
        setNotifType('error');
      }
    } catch {
      setNotif('Terjadi error jaringan!');
      setNotifType('error');
    }
  };

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();
    if (!form.email) {
      setNotif('Email wajib diisi!');
      setNotifType('error');
      return;
    }

    try {
      let res;
      if (isEdit) {
        res = await fetch(`${API_URL}/${form.id}`, {
          method: 'PUT',
          headers: { 
            'Content-Type': 'application/json', 
            Authorization: 'Bearer ' + token 
          },
          body: JSON.stringify(form)
        });
      } else {
        res = await fetch(API_URL, {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json', 
            Authorization: 'Bearer ' + token 
          },
          body: JSON.stringify(form)
        });
      }

      if (res.ok) {
        setNotif(isEdit ? 'User berhasil diupdate!' : 'User berhasil ditambah!');
        setNotifType('success');
        fetchUsers(token);
        handleCancel();
      } else {
        setNotif('Gagal menyimpan user!');
        setNotifType('error');
      }
    } catch {
      setNotif('Terjadi error jaringan!');
      setNotifType('error');
    }
  };

  const handleEdit = user => {
    setForm(user);
    setIsEdit(true);
  };

  const handleCancel = () => {
    setForm({ id: '', email: '', role: 'user' });
    setIsEdit(false);
  };

  return (
    <div className="user-management">
      <h2>Manajemen User</h2>
      <h3>{isEdit ? 'Edit User' : 'Tambah User Baru'}</h3>
      
      {notif && (
        <div style={{
          color: notifType === 'success' ? 'green' : 'red',
          marginBottom: 10
        }}>
          {notif}
        </div>
      )}

      <form onSubmit={handleSubmit} style={{ marginBottom: 20 }}>
        <div style={{ marginBottom: 10 }}>
          <input
            type="email"
            name="email"
            placeholder="Email"
            value={form.email}
            onChange={handleChange}
            required
            style={{ marginRight: 10 }}
          />
          
          <select
            name="role"
            value={form.role}
            onChange={handleChange}
            style={{ marginRight: 10 }}
          >
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>

          <button type="submit">
            {isEdit ? 'Update' : 'Tambah'}
          </button>
          
          {isEdit && (
            <button type="button" onClick={handleCancel} style={{ marginLeft: 10 }}>
              Batal
            </button>
          )}
        </div>
      </form>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ border: '1px solid #ddd', padding: 8 }}>ID</th>
              <th style={{ border: '1px solid #ddd', padding: 8 }}>Email</th>
              <th style={{ border: '1px solid #ddd', padding: 8 }}>Role</th>
              <th style={{ border: '1px solid #ddd', padding: 8 }}>Aksi</th>
            </tr>
          </thead>
          <tbody>
            {users.map(user => (
              <tr key={user.id}>
                <td style={{ border: '1px solid #ddd', padding: 8 }}>{user.id}</td>
                <td style={{ border: '1px solid #ddd', padding: 8 }}>{user.email}</td>
                <td style={{ border: '1px solid #ddd', padding: 8 }}>{user.role}</td>
                <td style={{ border: '1px solid #ddd', padding: 8 }}>
                  <button onClick={() => handleEdit(user)} style={{ marginRight: 5 }}>
                    Edit
                  </button>
                  <button onClick={() => handleDelete(user.id)}>
                    Hapus
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default UserManagement; 