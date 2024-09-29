// columns.tsx: contains column definitions for table components
"use client"

import React from 'react'
import { ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"
import { Link } from "react-router-dom"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import { ArrowUpDown } from "lucide-react"


// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Brother = {
    brotherID: string
    rollCall: number 
    firstName: string
    lastName: string
    status: string
    className: string
    email: string
    phoneNumber: string
}

export const brothersTableColumns: ColumnDef<Brother>[] = [
  {
    accessorKey: "rollCall",
    header: "Roll Call",
  },
  {
    accessorKey: "firstName",
    //header: "First Name",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          First Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      )
    },
  },
  {
    accessorKey: "lastName",
    header: "Last Name",
  },
    {
    accessorKey: "status",
    header: "Status",
  },
  {
    accessorKey: "className",
    header: "Class Name",
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "phoneNumber",
    header: "Phone Number",
  },
  {
      id: "actions",
    cell: ({ row }) => {
      const brother = row.original
 
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>

            <DropdownMenuItem onClick={() => navigator.clipboard.writeText(brother.brotherID)} >
              Copy Brother ID
            </DropdownMenuItem>
            <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(brother.firstName + " " + brother.lastName)} >
                Copy Full Name
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem>
                View Brother
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
    }
]

export type Event = {
    eventID: string
    eventName: string
    categoryName: string
    eventLocation: string
    eventDate: string
}

export const eventsTableColumns: ColumnDef<Event>[] = [
    {
        accessorKey: "eventID",
         header: "Event ID",
    },
    {
        accessorKey: "eventName",
        header: "Event Name",
        cell: ({ row }) => {
            const event = row.original
            return (
                <Link to={`/events/${event.eventID}/attendance`}>
                    <Button variant="ghost">
                        {event.eventName}
                    </Button>
                </Link>
            )
        },
    },
    {
        accessorKey: "categoryName",
        header: "Chair",
    },
    {
        accessorKey: "eventLocation",
        header: "Event Location",
    },
    {
        accessorKey: "eventDate",
        header: "Event Date",
    },
  {
      id: "actions",
    cell: ({ row }) => {
      const event = row.original
 
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>Actions</DropdownMenuLabel>

            <DropdownMenuItem onClick={() => navigator.clipboard.writeText(event.eventID)} >
              Copy Event ID
            </DropdownMenuItem>
            <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(event.eventName)} >
                Copy Event Name
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem>
                <Link to={`events/${event.eventID}/attendance`}>
                    View Event Attendance
                </Link>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
    }
]


export type EventAttendance = {
    brotherID: number
    firstName: string
    lastName: string
    rollCall: number
    attendanceStatus: string
}

export const eventAttendanceTableColumns: ColumnDef<EventAttendance>[] = [
    {
        accessorKey: "firstName",
         header: "First Name",
    },
    {
        accessorKey: "lastName",
        header: "Last Name",
    },
    {
        accessorKey: "rollCall",
        header: "Roll Call",
    },
    {
        accessorKey: "attendanceStatus",
        header: "Status",
    },
    {
        id: "actions",
        cell: ({ row }) => {
        const attendance = row.original

        return (
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>

                    <DropdownMenuItem onClick={() => navigator.clipboard.writeText(String(attendance.brotherID))} >
                        Copy Brother ID
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(attendance.firstName + " " + attendance.lastName)} >
                        Copy Full Name
                    </DropdownMenuItem>

                    <DropdownMenuSeparator />

                    <DropdownMenuItem>
                        View Brother
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
            )
        },
    }
]
