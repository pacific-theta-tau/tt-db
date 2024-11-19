// MainLayout.tsx: Main Layout to be used for all pages except for login page
import React from 'react';
import { Outlet } from 'react-router-dom'; // For nested routes

import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"


const MainLayout: React.FC = () => {
  return (
    <SidebarProvider>
        <AppSidebar />
        <SidebarInset>
            <SidebarTrigger />
            <div className="flex flex-1 flex-col gap-4 p-4">
                {/* This renders the component associated with the current route */}
                <Outlet />
            </div>
        </SidebarInset>
    </SidebarProvider>
  );
};

export default MainLayout;
