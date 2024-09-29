import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';

import LoginPage from './pages/Login'
import HomePage from './pages/Home'
import BrothersPage from './pages/Brothers'
import EventsPage from './pages/Events'
import EventAttendancePage from './pages/EventAttendance'
import NotFoundPage from './pages/NotFound'
import NavBar from './components/navbar'
import './App.css'


const App: React.FC = () => {
    return (
        <BrowserRouter>
            <div className="flex">
                <div>
                    <NavBar />
                </div>
                <div className="flex-grow p-8">
                    <Routes>
                        <Route path="/login" element={<LoginPage />} />
                        <Route path="/" element={<HomePage />} />
                        <Route path="/brothers" element={<BrothersPage />} />
                        <Route path="/events" element={<EventsPage />} />
                        <Route path="/events/:eventID/attendance" element={<EventAttendancePage />} />
                        <Route path="/404" element={<NotFoundPage />} />
                        <Route path="*" element={<Navigate replace to ="/404" />} />
                    </Routes>
                </div>
            </div>
        </BrowserRouter>
    );
}

export default App
