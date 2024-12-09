// edit-brothers-form.tsx: "Edit Row" form to be used by <SideRowSheet /> component
"use client"

import React from 'react'
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
import { ApiResponse, request } from '@/api/api'

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

import { Button } from "@/components/ui/button"

import { Input } from "@/components/ui/input"
import { Brother } from '@/components/columns'


const majors: readonly [string, ...string[]] = [
    'Civil Engineering',
    'Bioengineering',
    'Computer Engineering',
    'Computer Science',
    'Electrical Engineering',
    'Engineering Physics',
    'Mechanical Engineering',
]

const statuses: readonly [string, ...string[]] = [
    'Active',
    'Pre-Alumnus',
    'Alumnus',
    'Co-op',
    'Transferred',
    'Expelled',
]

const formSchema = z.object({
    firstName: z.string().min(1, "You must provide a First Name"),
    lastName: z.string().min(1, "You must provide a Last Name"),
    major: z.string().min(1, "You must provide a Major"),
    rollCall: z.number().min(1, "You must provide a Roll Call"),
    status: z.enum(statuses, {
                required_error: "You need to select status.",
            }),
    className: z.string().optional(),
    email: z.string().optional(),
    phoneNumber: z.string().optional(),
})


async function sendPatchRequest(data: z.infer<typeof formSchema>, brotherID: string): Promise<ApiResponse<Brother>>{
    /**
     * Mutation function to send patch request to modify row
     *
     * @param data - Data being sent as request body
     * @param brotherID - ID of the record to be changed
     * @returns Promise of an ApiResponse
     */
    const endpoint = `http://localhost:8080/api/brothers/${brotherID}`
    let result: ApiResponse<Brother>
    result = await request(endpoint, "PATCH", data)
    /* uncomment line below to test skeleton during loading */
    // await new Promise(f => setTimeout(f, 3000));
    return result
}

export function EditBrotherForm({rowData}: {rowData: Brother} ) {
  const { toast } = useToast()
  const brotherID = rowData.brotherID
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
        firstName: rowData.firstName,
        lastName: rowData.lastName,
        major: rowData.major,
        rollCall: rowData.rollCall,
        status: rowData.status,
        className: rowData.className,
        email: rowData.className,
        phoneNumber: rowData.phoneNumber,
    },
  })

  const queryClient = useQueryClient();
  // React Query mutation hook
  const mutation = useMutation(
  {
    mutationFn: (data: z.infer<typeof formSchema>) => sendPatchRequest(data, brotherID),
    onSuccess: (data) => {
        // TODO: use "message" field for toast description
        toast({
            title: "Success!",
            description: "Record updated successfully.",
        })
      // Invalidate table data query to reload the table
      queryClient.invalidateQueries({ queryKey: ["brothersTableData"] });
    },
    onError: (error) => {
        // Make toast destructive
        toast({
            title: "Uh oh! Something went wrong.",
            description: "Failed to update record.",
            variant: "destructive",
            //action: <ToastAction></ToastAction>,
        })
    }
  });

  async function onSubmit(data: z.infer<typeof formSchema>) {
      mutation.mutate(data)
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="firstName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>First Name *</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormDescription>
                {}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="lastName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Last Name *</FormLabel>
              <FormControl>
                <Input placeholder={rowData.lastName} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="major"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Major *</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                  <FormControl>
                          <SelectTrigger className="w-[250px]">
                              <SelectValue placeholder="Select Major" />
                          </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                        {majors.map((major, index) => (
                          <SelectItem value={major} key={index}>{major}</SelectItem>
                        ))}
                  </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="rollCall"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Roll Call *</FormLabel>
              <FormControl>
                <Input
                    type="number" {...field}
                    onChange={event => field.onChange(+event.target.value)}
                />
              </FormControl>
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
 
         <FormField
          control={form.control}
          name="className"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Class</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
 
         <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

 
         <FormField
          control={form.control}
          name="phoneNumber"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Phone Number</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />



        <Button type="submit">Submit</Button>
      </form>
    </Form>
  )
}

