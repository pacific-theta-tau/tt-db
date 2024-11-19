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



/*
 * Helper function to get current Semester + Year in URL format
 * */
export function getSeasonYear(): string {
  const today = new Date();
  const month = today.getMonth(); // Months are 0-indexed: January is 0, December is 11
  const year = today.getFullYear();

  // Determine the season based on the month
  const season = month < 6 ? 'Spring' : 'Fall';

  return encodeURIComponent(`${season} ${year}`);
}

export function AppSidebar() {
    const semester = getSeasonYear()
    console.log("semester:", semester)

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
        url: `/actives/${semester ? semester : "Spring%202024"}`,
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
        <p className="font-bold pl-2">Pacific Theta Tau Database</p>
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

            { /* dynamically get the url for actives of current semester */ }
                <SidebarMenuItem className="mx-2 mb-1">
                  <SidebarMenuButton asChild>
                    <Link to={`/actives/${semester}`}>
                      <Users />
                      <span>Actives</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>


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

