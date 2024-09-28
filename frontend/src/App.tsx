import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
//import BrothersTable from './components/brothers-table'
//import EventsTable from './components/events-table'
import HomePage from './pages/Home'
import BrothersPage from './pages/Brothers'
import EventsPage from './pages/Events'
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
                        <Route path="/" element={<HomePage />} />
                        <Route path="/brothers" element={<BrothersPage />} />
                        <Route path="/events" element={<EventsPage />} />
                    </Routes>
                </div>
            </div>
        </BrowserRouter>
    );
}

export default App
