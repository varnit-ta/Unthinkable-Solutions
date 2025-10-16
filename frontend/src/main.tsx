/**
 * Application Entry Point
 * 
 * This file initializes the React application with necessary providers:
 * - React.StrictMode for development checks
 * - BrowserRouter for client-side routing
 * - AuthProvider for authentication state management
 * - Toaster for notification system
 * 
 * @module main
 */

import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App.tsx'
import './index.css'
import { AuthProvider } from './auth.tsx'
import { Toaster } from './components/ui/sonner'

/**
 * Initialize and render the React application
 */
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider>
        <App />
        <Toaster />
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>,
)
