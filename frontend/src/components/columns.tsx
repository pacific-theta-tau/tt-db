// columns.tsx: contains column definitions for table components
"use client"

import React from 'react'
import { ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal, Clipboard, Pencil, Trash2 } from "lucide-react"
import { Link } from "react-router-dom"
import { ArrowUpDown } from "lucide-react"
import { Button, buttonVariants } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { DeleteAlertDialog } from '@/components/delete-alert-dialog'
import SideRowSheet from './sheet/side-row-sheet'
import { EditBrotherForm } from './sheet/forms/edit-brothers-form'
import { EditEventsForm } from './sheet/forms/edit-events-form'
import { EditEventAttendanceForm } from './sheet/forms/edit-attendance-form'
import { eventsQueryKey } from '@/components/events-table'
import { attendanceQueryKey } from '@/pages/EventAttendance'
import { activesQueryKey } from '@/pages/Actives'


// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Brother = {
    brotherID: string
    rollCall: number 
    firstName: string
    lastName: string
    major: string
    status: string
    className: string
    email: string
    phoneNumber: string
}

export const brothersTableColumns: ColumnDef<Brother>[] = [
    {
        accessorKey: "rollCall",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    Roll Call
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
    },
    {
        accessorKey: "firstName",
        //header: "First Name",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    First Name
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
    },
    {
    accessorKey: "lastName",
         header: ({ column }) => {
             return (
                 <Button
                     variant="ghost"
                     onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                     className="px-0"
                 >
                     Last Name 
                 <ArrowUpDown className="ml-2 h-4 w-4" />
                 </Button>
                )
         },
    },
    {
        accessorKey: "major",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    Major 
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
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
            const deleteEndpoint = "/api/brothers"
            const deleteBodyParams = {
                "rollCall": brother.rollCall
            }
 
            return (
                // Setting modal=false to fix `AlertDialog` side-effects of not letting click anything
                <DropdownMenu modal={false}>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <span className="sr-only">Open menu</span>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Actions</DropdownMenuLabel>

                        <SideRowSheet
                            title="Edit Row"
                            description=""
                            FormType={
                                <EditBrotherForm 
                                    rowData={brother}
                                />
                            }
                            trigger={
                                <DropdownMenuItem onSelect={ (e) => e.preventDefault() } >
                                    <Pencil className="mr-2 h-4 w-4"/> Edit
                                </DropdownMenuItem>
                            }
                        />

                        <DeleteAlertDialog
                            endpoint={ deleteEndpoint }
                            body={ deleteBodyParams }
                            trigger={
                                <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  <span>Delete</span>
                                </DropdownMenuItem>
                            }
                            queryKey="brothersTableData"
                            />

                        <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(brother.firstName + " " + brother.lastName)} >
                             <Clipboard className="mr-2 h-4 w-4"/> Copy Full Name
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
        accessorKey: "eventName",
        header: ({ column }) => {
            return (
                 <Button
                     variant="ghost"
                     onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                     className="px-0"
                 >
                    Event Name 
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                 </Button>
            )
        },

        cell: ({ row }) => {
            const event = row.original
            return (
                <Link to={`/events/${event.eventID}/attendance`}>
                    <Button variant="link" className="px-0">
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
        header: ({ column }) => {
            return (
                <Button
                        variant="ghost"
                        onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                        className="px-0"
                >
                    Event Date 
                    <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
            )
        },
    },
    {
        id: "actions",
        cell: ({ row }) => {
            const event = row.original
            const deleteEndpoint = "/api/events"
            const deleteBodyParams = {
                "eventID": parseInt(event.eventID)
            }
     
            return (
                // Setting modal=false to fix `AlertDialog` side-effects of not letting click anything
                <DropdownMenu modal={false}>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                            <span className="sr-only">Open menu</span>
                            <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>

                    <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Actions</DropdownMenuLabel>

                        <SideRowSheet
                            title="Edit Row"
                            description=""
                            FormType={
                                <EditEventsForm
                                    rowData={event}
                                />
                            }
                            trigger={
                                <DropdownMenuItem onSelect={ (e) => e.preventDefault() } >
                                    <Pencil className="mr-2 h-4 w-4"/> Edit
                                </DropdownMenuItem>
                            }
                        />
                                                
                        <DeleteAlertDialog
                            endpoint={ deleteEndpoint }
                            body={ deleteBodyParams }
                            trigger={
                                <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  <span>Delete</span>
                                </DropdownMenuItem>
                            }
                            queryKey={ eventsQueryKey }
                            >
                        </DeleteAlertDialog>

                        <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(event.eventName)} >
                            <Clipboard className="mr-2 h-4 w-4" /> Copy Event Name
                        </DropdownMenuItem>

                        <DropdownMenuSeparator />

                        <DropdownMenuItem>
                            <Link to={`/events/${event.eventID}/attendance`}>
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
    eventID?: number
    firstName: string
    lastName: string
    rollCall: number
    attendanceStatus: string
}

export const eventAttendanceTableColumns: ColumnDef<EventAttendance>[] = [
    {
        accessorKey: "firstName",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    First Name
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
           )
        },
    },
    {
        accessorKey: "lastName",
        header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
              className="px-0"
            >
                Last Name 
              <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
    },
    {
        accessorKey: "rollCall",
        header: ({ column }) => {
          return (
            <Button
              variant="ghost"
              onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
              className="px-0"
            >
                Roll Call 
            <ArrowUpDown className="ml-2 h-4 w-4" />
            </Button>
          )
        },
    },
    {
        accessorKey: "attendanceStatus",
        header: "Status",
    },
    {
        id: "actions",
        cell: ({ row }) => {
        const attendance = row.original
        const deleteEndpoint = "/api/attendance"
        const deleteBodyParams = {
            "brotherID": attendance.brotherID,
            "eventID": attendance.eventID
        }

        return (
            <DropdownMenu modal={false}>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>

                    <SideRowSheet
                            title="Edit Row"
                            description=""
                            FormType={
                                <EditEventAttendanceForm
                                    rowData={attendance}
                                />
                            }
                            trigger={
                                <DropdownMenuItem onSelect={ (e) => e.preventDefault() } >
                                    <Pencil className="mr-2 h-4 w-4"/> Edit
                                </DropdownMenuItem>
                            }
                    />

                    <DeleteAlertDialog
                            endpoint={ deleteEndpoint }
                            body={ deleteBodyParams }
                            trigger={
                                <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  <span>Delete</span>
                                </DropdownMenuItem>
                            }
                            queryKey={attendanceQueryKey}
                            >
                    </DeleteAlertDialog>

                    <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(attendance.firstName + " " + attendance.lastName)} >
                        <Clipboard className="mr-2 h-4 w-4" /> Copy Full Name
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

// Used for search table in Attendance Table's "add row" sheet
export const rollCallSearchColumns: ColumnDef<Brother>[] = [
    {
        accessorKey: 'rollCall',
        header: 'Roll Call',
    },
    {
        accessorKey: 'firstName',
         header: 'First Name',
    },
    {
        accessorKey: 'lastName',
        header: 'Last Name',
    },
    {
        accessorKey: 'class',
        header: 'Class',
    }
]


export type BrotherStatus = {
    brotherID: number,
    rollCall: number,
    firstName: string,
    lastName: string,
    major: string,
    class: string,
    status: string,
    semesterID: number,
    semesterLabel: string,
}

export const brotherStatusTableColumns: ColumnDef<BrotherStatus>[] = [
     {
        accessorKey: "rollCall",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    Roll Call
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
    },
    {
        accessorKey: "firstName",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    First Name
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
    },
    {
        accessorKey: "lastName",
         header: ({ column }) => {
             return (
                 <Button
                     variant="ghost"
                     onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                     className="px-0"
                 >
                     Last Name 
                 <ArrowUpDown className="ml-2 h-4 w-4" />
                 </Button>
                )
         },
    },
    {
        accessorKey: "major",
        header: ({ column }) => {
            return (
                <Button
                    variant="ghost"
                    onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
                    className="px-0"
                >
                    Major 
                <ArrowUpDown className="ml-2 h-4 w-4" />
                </Button>
               )
        },
    },
    {
        accessorKey: "status",
        header: "Status",
    },
    {
        accessorKey: "class",
        header: "Class Name",
    },
    {
          id: "actions",
        cell: ({ row }) => {
            const brotherStatus = row.original
            const brotherID = brotherStatus.brotherID
            const semesterID = brotherStatus.semesterID
            const deleteEndpoint = `/v1/brothers/${brotherID}/statuses/${semesterID}`
     
          return (
            // Setting modal=false to fix `AlertDialog` side-effects of not letting click anything
            <DropdownMenu modal={false}>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                  <span className="sr-only">Open menu</span>
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>Actions</DropdownMenuLabel>

                <DropdownMenuItem onClick={ () => console.log("Edit row") } >
                    <Pencil className="mr-2 h-4 w-4"/>Edit
                </DropdownMenuItem>

                <DeleteAlertDialog
                        endpoint={ deleteEndpoint }
                        trigger={
                            <DropdownMenuItem onSelect={(e) => e.preventDefault()}>
                              <Trash2 className="mr-2 h-4 w-4" />
                              <span>Delete</span>
                            </DropdownMenuItem>
                        }
                        queryKey="activesTableData"
                        >
                </DeleteAlertDialog>

                <DropdownMenuItem onClick={ () => navigator.clipboard.writeText(brotherStatus.firstName + " " + brotherStatus.lastName)} >
                     <Clipboard className="mr-2 h-4 w-4"/>Copy Full Name
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
