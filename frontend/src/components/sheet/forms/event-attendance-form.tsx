// form-types.tsx: Defines form components for different data table pages
// These forms should be used as props to `<AddSheet />`
"use client"

import React, { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom';
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
import { Search } from "lucide-react"
import { useReactTable, getCoreRowModel, getFilteredRowModel, flexRender, ColumnDef } from '@tanstack/react-table'
import { rollCallSearchColumns } from '@/components/columns';
import { Brother } from "@/components/columns"


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

export function EventAttendanceForm() {
    const [rollCall, setRollCall] = useState(0)
    const [globalFilter, setGlobalFilter] = useState("")
    const [searchData, setSearchData] = useState<Brother[]>([])
    const [isDialogOpen, setIsDialogOpen] = useState(false)

    const { toast } = useToast()
    const { eventID } = useParams<{ eventID: string }>();

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
        },
    })

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/brothers"
        const fetchData = async () => {
             try {
                const response = await fetch(endpoint, {
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                if (!response.ok) {
                  throw new Error('Network response was not ok');
                }
                const result: Brother[] = await response.json();
                console.log('result:', result)
                setSearchData(result);
            } catch (error) {
                console.log('Error fetching data:', error);
                throw error;
            } finally {
                /* uncomment line below to test skeleton during loading */
                // await new Promise(f => setTimeout(f, 3000));
            }
        }
        fetchData()
       }, []);


    const table = useReactTable({
        data: searchData,
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

  async function onSubmit(data: z.infer<typeof formSchema>) {
    const endpoint = "http://localhost:8080/api/events" + eventID + "/attendance"
    let result: any
    const body = {
            "eventID": eventID,
            "brotherID": "",
            "rollCall": rollCall,
            "status": data.status,
    }
    try {
        
        const response = await fetch(endpoint, {
            method: 'POST',
            body: JSON.stringify(body),
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json',
            }
        });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        result  = await response.json();
        console.log('result:', result)
    } catch (error) {
        console.log('Error fetching data:', error);
        throw error;
    } finally {
        /* uncomment line below to test skeleton during loading */
        // await new Promise(f => setTimeout(f, 3000));
        console.log(result)
        toast({
            title: "You submitted the following values:",
            description: (
                <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
                    <code className="text-white">{JSON.stringify(body, null, 2)}</code>
                </pre>
            ),
        })
    }
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
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}

