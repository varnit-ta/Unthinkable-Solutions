import React, { useState } from 'react'
import { api } from '../api'
import { useAuth } from '../auth'

export default function RegisterPage() {
  const { setToken } = useAuth()
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [err, setErr] = useState<string | null>(null)

  async function submit(e: React.FormEvent) {
    e.preventDefault()
    try {
      const res = await api.register(username, email, password)
      setToken(res.token)
    } catch (e: any) {
      setErr(e.message || 'register failed')
    }
  }

  return (
    <div style={{ padding: 16 }}>
      <h2>Register</h2>
      <form onSubmit={submit}>
        <div>
          <label>Username</label>
          <input value={username} onChange={(e) => setUsername(e.target.value)} />
        </div>
        <div>
          <label>Email</label>
          <input value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <div>
          <label>Password</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <button type="submit">Register</button>
      </form>
      {err && <div style={{ color: 'red' }}>{err}</div>}
    </div>
  )
}
