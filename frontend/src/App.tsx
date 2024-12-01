import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"

import LoginPage from './pages/Login'
import HomePage from './pages/Home'
import BrothersPage from './pages/Brothers'
import ActivesPage from './pages/Actives'
import EventsPage from './pages/Events'
import EventAttendancePage from './pages/EventAttendance'
import NotFoundPage from './pages/NotFound'
import { NavBar2 } from './components/navbar'
import './App.css'

import { ThemeProvider } from "@/components/theme-provider"
import { Toaster } from "@/components/ui/toaster"
import MainLayout from './layouts/MainLayout';


const App: React.FC = () => {
    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <Toaster />
            <BrowserRouter>
                <Routes>
                    <Route element={<MainLayout />} >
                        <Route path="/" element={<HomePage />} />
                        <Route path="/brothers" element={<BrothersPage />} />
                        <Route path="/actives" element={<ActivesPage />} />
                        <Route path="/events" element={<EventsPage />} />
                        <Route path="/events/:eventID/attendance" element={<EventAttendancePage />} />
                    </Route>
                    
                    { /* Routes without navbar */ }
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/404" element={<NotFoundPage />} />
                    <Route path="*" element={<Navigate replace to ="/404" />} />
                </Routes>
            </BrowserRouter>
        </ThemeProvider>
    );
}

export default App
