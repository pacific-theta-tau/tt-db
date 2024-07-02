"use client"

import { ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Brother = {
    id: string
    rollCall: number 
    firstName: string
    lastName: string
    status: string
    className: string
    email: string
    phoneNumber: string
}

export const columns: ColumnDef<Brother>[] = [
  {

    accessorKey: "rollCall",
    header: "Roll Call",
  },
  {

    accessorKey: "firstName",
    header: "First Name",
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
      const payment = row.original
 
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
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(payment.id)}
            >
              Copy payment ID
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>View customer</DropdownMenuItem>
            <DropdownMenuItem>View payment details</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
    }
]

