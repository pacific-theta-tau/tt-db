// app-sidebar.tsx: new navbar component. `navbar.tsx` is deprecated
import React from 'react';
import { Link, useLocation } from "react-router-dom"
import { Home, Table, Users, } from "lucide-react"
import ThemeToggle from './theme-toggle'

import {
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"


export function AppSidebar() {
    // Menu items.
    let items = [
      {
        title: "Home",
        url: "/",
        icon: Home,
      },
      {
        title: "Members",
        url: "/brothers",
        icon: Table,
      },
      {
        title: "Actives",
        url: `/actives`,
        icon: Users,
      },
      {
        title: "Events",
        url: "/events",
        icon: Table,
      },
    ]

  return (
    <Sidebar>
      <SidebarHeader>
        <p className="font-bold pl-2">Theta Tau DB</p>
      </SidebarHeader >
      <SidebarContent>
        <SidebarGroup />
        <SidebarGroupLabel className="pl-2">Platform</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map((item) => (
                <SidebarMenuItem className="mx-2 mb-1" key={item.title}>
                  <SidebarMenuButton asChild>
                    <Link to={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}

              <SidebarMenuItem className="mx-2 mb-1">
                  <ThemeToggle />
              </ SidebarMenuItem>
               
            </SidebarMenu>
          </SidebarGroupContent>
        <SidebarGroup />
      </SidebarContent>
    </Sidebar>
  )
}

