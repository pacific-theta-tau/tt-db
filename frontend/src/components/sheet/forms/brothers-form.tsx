// form-types.tsx: Defines form components for different data table pages
// These forms should be used as props to `<AddSheet />`
"use client"

import React from 'react'
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
import { ApiResponse, requestPOST } from '@/api/api'

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
    firstName: z.string({
        required_error: "You must provide a first name"
    }),
    lastName: z.string({
        required_error: "You must provide a last name"
    }),
    major: z.string({
        required_error: "You must provide a major",
    }),
    rollCall: z.number({
        required_error: "You must provide a roll call"
    }),
    status: z.enum(statuses, {
                required_error: "You need to select status.",
            }),
    className: z.string().optional(),
    email: z.string().optional(),
    phoneNumber: z.string().optional(),
})

export function BrotherForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      phoneNumber: "",
    },
  })

  const { toast } = useToast()
  async function onSubmit(data: z.infer<typeof formSchema>) {
    const endpoint = "http://localhost:8080/api/brothers"
    let result: ApiResponse<Brother>
    try {
        const body = JSON.stringify(data)
        result = await requestPOST(endpoint, body)
        console.log('result:', result)
    } catch (error) {
        console.log('Error fetching data:', error);
        throw error;
    } finally {
        /* uncomment line below to test skeleton during loading */
        // await new Promise(f => setTimeout(f, 3000));
        toast({
            title: "You submitted the following values:",
            description: (
                <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
                    <code className="text-white">{JSON.stringify(data, null, 2)}</code>
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
                <Input placeholder="" {...field} />
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
                          <SelectTrigger className="w-[180px]">
                              <SelectValue placeholder="Select Major" />
                          </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                        {majors.map((major) => (
                          <SelectItem value={major}>{major}</SelectItem>
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

