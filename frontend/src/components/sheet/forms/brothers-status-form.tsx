// form-types.tsx: Defines form components for different data table pages
// These forms should be used as props to `<AddSheet />`
"use client"

import React, { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams } from 'react-router-dom';
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
import { Search } from "lucide-react"
import { useReactTable, getCoreRowModel, getFilteredRowModel, flexRender, ColumnDef } from '@tanstack/react-table'
import { rollCallSearchColumns } from '@/components/columns';
import { Brother, BrotherStatus } from "@/components/columns"
import { request, ApiResponse } from "@/api/api"
import { activesQueryKey } from "@/pages/Actives"


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
    'Active',
    'Pre-Alumnus',
    'Alumnus',
    'Co-op',
    'Transferred',
    'Expelled',
]

const formSchema = z.object({
        rollCall: z.number({
            required_error: "You must provide a roll call"
        }),
        semester: z.string({
            required_error: "You must provide a semester"
        }),
        status: z.enum(statuses, {
            required_error: "You need to select status.",
        }),
    })


async function fetchSearchData() {
    console.log("CALLED fetchSearchData")
    const endpoint = "http://localhost:8080/api/brothers"
    const responseSearch: ApiResponse<Brother[]> = await request(endpoint, 'GET')
    console.log('responseSearch:', responseSearch)

    return responseSearch.data
}


async function fetchSemestersData() {
    console.log("CALLED fetchSemestersData")
    const endpoint2 = "http://localhost:8080/api/semesters"
    const responseSemesters: ApiResponse<string[]> = await request(endpoint2, 'GET')
    console.log('responseSemesters:', responseSemesters)
    return responseSemesters.data
}


async function sendPostRequest(data: z.infer<typeof formSchema>, semester: string, brotherID: string) {
    const endpoint = `http://localhost:8080/api/semesters/${semester}/statuses`
    const body = {
            "brotherID": parseInt(brotherID),
            "status": data.status,
    }
    const result: ApiResponse<BrotherStatus[]> = await request(endpoint, 'POST', body)
    return result.data
}


export function BrotherStatusForm({selectedSemester}: { selectedSemester: string }){
    console.log("Selected Semester:", selectedSemester)
    const { toast } = useToast()
    const [rollCall, setRollCall] = useState(0)
    const [brotherID, setBrotherID] = useState("")
    const [globalFilter, setGlobalFilter] = useState("")
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            semester: selectedSemester
        },
    })

    // React Query and Mutation Hooks
    // TODO: handle isLoading and isError states for search dialog and semester dropdown
    // searchData is for "Search Brother" search field
    const { data: searchData, isLoading, isError } = useQuery({queryKey: ["brotherSearchData"], queryFn: fetchSearchData})
    // semestersData is for the "Select Semester" dropdown field
    const { data: semestersData } = useQuery({queryKey: ["semestersData"], queryFn: fetchSemestersData})
    const queryClient = useQueryClient();
    const mutation = useMutation(
    {
        mutationFn: (data: z.infer<typeof formSchema>) => sendPostRequest(data, selectedSemester, brotherID),
        onSuccess: (data) => {
            // TODO: use "message" field for toast description
            toast({
                title: "You submitted the following values:",
                description: (
                    <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
                        <code className="text-white">{JSON.stringify(data, null, 2)}</code>
                    </pre>
                ),
            })
              // Invalidate table data query to reload the table
              queryClient.invalidateQueries({ queryKey: [activesQueryKey] });
            },
        onError: (error) => {
            // Make toast destructive
            toast({
                title: "Failed to submit ",
                variant: "destructive",
                //action: <ToastAction></ToastAction>,
                description: (
                    <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
                        <code className="text-white">{error.message}</code>
                    </pre>
                ),
            })
        }
    })

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
        setBrotherID(rollCallSearch.brotherID)
        form.setValue("rollCall", rollCallSearch.rollCall)
        setIsDialogOpen(false)
    }

    async function onSubmit(data: z.infer<typeof formSchema>) {
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
                          placeholder="Enter member Roll Call"
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
                          <DialogContent className="sm:max-w-[500px] overflow-y-scroll sm:max-h-[700px]">
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
            name="semester"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Semester *</FormLabel>
                <Select onValueChange={field.onChange} defaultValue={selectedSemester}>
                    <FormControl>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Select Semester" />
                            </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                          {semestersData && semestersData.length > 0 ? (
                              semestersData.map((semester, index) => (
                                  <SelectItem key={index.toString()} value={semester}>{semester}</SelectItem>
                              ))
                          ) : (
                              <SelectItem value={selectedSemester}>Loading...</SelectItem>
                          )
                          }
                    </SelectContent>
                </Select>
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
                <Select onValueChange={field.onChange} defaultValue="">
                    <FormControl>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Select Status" />
                            </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                          {statuses.map((status, index) => (
                            <SelectItem value={status} key={index.toString()}>{status}</SelectItem>
                          ))}
                    </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          {/* TODO: add `disabled` attribute for useQuery isLoading? */}
          <Button type="submit">Submit</Button>
        </form>
      </Form>
    )
}

