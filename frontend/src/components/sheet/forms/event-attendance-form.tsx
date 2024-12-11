// form-types.tsx: Defines form components for different data table pages
// These forms should be used as props to `<AddSheet />`
"use client"

import React, { useState, useEffect } from 'react'
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query'
import { useParams } from 'react-router-dom';
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
import { Search } from "lucide-react"
import { useReactTable, getCoreRowModel, getFilteredRowModel, flexRender, ColumnDef } from '@tanstack/react-table'
import { rollCallSearchColumns } from '@/components/columns';
import { Brother, EventAttendance } from "@/components/columns"
import { request, ApiResponse } from '@/api/api';
import { attendanceQueryKey } from '@/pages/EventAttendance'


// UI imports
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"



const statuses: readonly [string, ...string[]] = [
    'Present',
    'Absent',
    'Excused',
]

const formSchema = z.object({
    rollCall: z.number({
        required_error: "You must provide a roll call"
    }),
    status: z.enum(statuses, {
        required_error: "You need to select status.",
    }),
})

async function fetchSearchData() {
    /**
    * Mutation function to create new event row from form data
    *
    * @param data - Form data to be sent in request body
    * @returns A Promise with Event data
    */
    const endpoint = "http://localhost:8080/api/brothers"
    const response: ApiResponse<Brother[]> = await request(endpoint, 'GET')
    return response.data
}

async function sendPostRequest(data: z.infer<typeof formSchema>, eventID: string, rollCall: number) {
    /**
    * Mutation function to create new attendance record from form data
    *
    * @param data - Form data to be sent in request body
    * @param eventID - ID of the event related to new attendance record
    * @param rollCall - rollCall of member
    * @returns A Promise with Event data
    */
    const endpoint = "http://localhost:8080/api/events/" + eventID + "/attendance"
    const body = {
            "eventID": eventID,
            "brotherID": "",
            "rollCall": rollCall,
            //"attendanceStatus": data.status,
            "attendanceStatus": data.status,
    }
    const result: ApiResponse<EventAttendance> = await request(endpoint, 'POST', body)
    return result
}

export function EventAttendanceForm() {
    const { toast } = useToast()
    const [rollCall, setRollCall] = useState(0)
    const [globalFilter, setGlobalFilter] = useState("")
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const { eventID = '0' } = useParams<{ eventID: string }>();
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
        },
    })

    // React Hook and table for member search field
    const { data: searchData, isLoading, isFetching, isError } = useQuery({queryKey: ["brotherSearchData"], queryFn: fetchSearchData })
    const table = useReactTable({
        data: searchData ?? [],
        columns: rollCallSearchColumns,
        getCoreRowModel: getCoreRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        onGlobalFilterChange: setGlobalFilter,
        state: {
            globalFilter,
        },
    })
    const handleSelectMember = (rollCallSearch: Brother) => {
        setRollCall(rollCallSearch.rollCall)
        form.setValue("rollCall", rollCallSearch.rollCall)
        setIsDialogOpen(false)
    }

    // React Hook and Mutation for Attendance Table
    const queryClient = useQueryClient();
    const mutation = useMutation(
    {
      mutationFn: (data: z.infer<typeof formSchema>) => sendPostRequest(data, eventID, rollCall),
      onSuccess: (data) => {
          // TODO: use "message" field for toast description
          toast({
              title: "Success!",
              description: "Added new attendance record to database.",
          })
        // Invalidate table data query to reload the table
        queryClient.invalidateQueries({ queryKey: [attendanceQueryKey] });
      },
      onError: (error) => {
          // Make toast destructive
          toast({
              title: "Uh oh! Something went wrong.",
              description: "Failed to add new attendance record.",
              variant: "destructive",
              //action: <ToastAction></ToastAction>,
          })
      }
    });
  
    async function onSubmit(data: z.infer<typeof formSchema>) {
      //const endpoint = "http://localhost:8080/api/attendance"
      mutation.mutate(data)
    }
  
    return (
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
                control={form.control}
                name="rollCall"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Roll Call</FormLabel>
                    <FormControl>
                      <div className="flex">
                        <Input
                          placeholder="Enter your Member ID"
                          {...field}
                          className="flex-grow"
                          onChange={
                              event => {
                                  field.onChange(+event.target.value)
                                  setRollCall(+event.target.value)
                              }
                          }/>
                        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                          <DialogTrigger asChild>
                            <Button variant="outline" className="ml-2">
                              <Search className="h-4 w-4" />
                            </Button>
                          </DialogTrigger>
                          <DialogContent className="sm:max-w-[500px]">
                            <DialogHeader>
                              <DialogTitle>Search Brothers</DialogTitle>
                            </DialogHeader>
                            <div className="py-4">
                              <Input
                                placeholder="Search by name..."
                                value={globalFilter ?? ""}
                                onChange={(e) => setGlobalFilter(String(e.target.value))}
                                className="mb-4"
                              />
                              <Table>
                                <TableHeader>
                                  {table.getHeaderGroups().map((headerGroup) => (
                                    <TableRow key={headerGroup.id}>
                                      {headerGroup.headers.map((header) => (
                                        <TableHead key={header.id}>
                                          {flexRender(
                                            header.column.columnDef.header,
                                            header.getContext()
                                          )}
                                        </TableHead>
                                      ))}
                                    </TableRow>
                                  ))}
                                </TableHeader>
                                <TableBody>
                                  {table.getRowModel().rows.map((row) => (
                                    <TableRow 
                                      key={row.id} 
                                      onClick={() => handleSelectMember(row.original)}
                                      className="cursor-pointer hover:bg-muted"
                                    >
                                      {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id}>
                                          {flexRender(
                                            cell.column.columnDef.cell,
                                            cell.getContext()
                                          )}
                                        </TableCell>
                                      ))}
                                    </TableRow>
                                  ))}
                                </TableBody>
                              </Table>
                            </div>
                          </DialogContent>
                        </Dialog>
                      </div>
                    </FormControl>
                    <FormDescription>
                      Enter Roll Call or search for Brother
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
  
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Status *</FormLabel>
                <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Select Status" />
                            </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                          {statuses.map((status) => (
                            <SelectItem value={status}>{status}</SelectItem>
                          ))}
                    </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" disabled={isLoading || isFetching || isError}>Submit</Button>
        </form>
      </Form>
    )
}

