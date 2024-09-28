import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
//import BrothersTable from './components/brothers-table'
//import EventsTable from './components/events-table'
import HomePage from './pages/Home'
import BrothersPage from './pages/Brothers.tsx'
import EventsPage from './pages/Events.tsx'
import './App.css'

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/brothers" element={<BrothersPage />} />
                <Route path="/events" element={<EventsPage />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App
