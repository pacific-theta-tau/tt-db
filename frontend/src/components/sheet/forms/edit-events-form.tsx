// form-types.tsx: Defines form components for different data table pages
// These forms should be used as props to `<AddSheet />`
"use client"

import React from 'react'
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useToast } from "@/hooks/use-toast"
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
import { format } from "date-fns"
import { Calendar as CalendarIcon } from "lucide-react"
import { cn } from "@/lib/utils"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

import { Event } from "@/components/columns"
import { eventsQueryKey } from "@/components/events-table"
import { ApiResponse, request } from '@/api/api'

// TODO: Add rest (or fill this dynamically with a db query for all categories)
const eventCategories: readonly [string, ...string[]] = [
    'Brotherhood',
    'Professional Development',
    'Community Service',
]

const formSchema = z.object({
    eventName: z.string({
        required_error: "You must provide a first name"
    }),
    categoryName: z.enum(eventCategories, {
                required_error: "You need to select status.",
            }),
    eventDate: z.date({
        required_error: "You must provide a date"
    }),
    eventLocation: z.string({ }).optional(),
})


async function sendPatchRequest(data: z.infer<typeof formSchema>, eventID: string): Promise<ApiResponse<Event[]>> {
    /**
    * Mutation function to create new event row from form data
    *
    * @param data - Form data to be sent in request body
    * @returns A Promise with Event data
    */
    const endpoint = `http://localhost:8080/api/events/${eventID}`
    const result: ApiResponse<Event[]> = await request(endpoint, 'PATCH', data);
    return result
}

interface EditEventsFormProps {
    rowData: Event
    onClose?: () => void;
}

export function EditEventsForm({ rowData, onClose } : EditEventsFormProps) {
    const { toast } = useToast()
    // Update display date with month/year selectors
    const [calendarDate, setCalendarDate] = React.useState(new Date())
    const form = useForm<z.infer<typeof formSchema>>({
      resolver: zodResolver(formSchema),
      defaultValues: {
          eventName: rowData.eventName,
          categoryName:rowData.categoryName,
          eventDate: new Date(rowData.eventDate),
          eventLocation: rowData.eventLocation,
      },
    })

    // React Query mutation hook
    const queryClient = useQueryClient();
    const mutation = useMutation(
    {
      mutationFn: (data: z.infer<typeof formSchema>) => sendPatchRequest(data, rowData.eventID),
      onSuccess: () => {
          // TODO: use "message" field for toast description
          toast({
              title: "Success!",
              description: "Added new event record to the database.",
          })
          // Invalidate table data query to reload the table
          queryClient.invalidateQueries({ queryKey: [eventsQueryKey] });
      },
      onError: (error) => {
          // Make toast destructive
          toast({
              title: "Uh oh! Something went wrong.",
              variant: "destructive",
              //action: <ToastAction></ToastAction>,
              description: `Failed to add new event record to database. Error: ${error.message}`,
          })
      },
      //onSettled: onClose,
    });

    async function onSubmit(data: z.infer<typeof formSchema>) {
        mutation.mutate(data)
    }

    const handleMonthChange = (month: string) => {
      const newDate = new Date(calendarDate)
      newDate.setMonth(parseInt(month))
      setCalendarDate(newDate)
    }

    const handleYearChange = (year: string) => {
      const newDate = new Date(calendarDate)
      newDate.setFullYear(parseInt(year))
      setCalendarDate(newDate)
    }

    const years = Array.from({ length: 11 }, (_, i) => {
      const year = new Date().getFullYear() - 5 + i
      return { value: year.toString(), label: year.toString() }
    })

    const months = [
      { value: "0", label: "January" },
      { value: "1", label: "February" },
      { value: "2", label: "March" },
      { value: "3", label: "April" },
      { value: "4", label: "May" },
      { value: "5", label: "June" },
      { value: "6", label: "July" },
      { value: "7", label: "August" },
      { value: "8", label: "September" },
      { value: "9", label: "October" },
      { value: "10", label: "November" },
      { value: "11", label: "December" },
    ]

    return (
      <Form {...form} >
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="eventName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Event Name *</FormLabel>
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
            name="categoryName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Event Category *</FormLabel>
                <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Select Status" />
                            </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                          {eventCategories.map((category) => (
                            <SelectItem value={category}>{category}</SelectItem>
                          ))}
                    </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="eventLocation"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Event Location</FormLabel>
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
            name="eventDate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Event Date</FormLabel>
                <Popover modal={true}>
                  <PopoverTrigger asChild>
                    <FormControl>
                      <Button
                        variant={"outline"}
                        className={cn(
                          "w-[280px] justify-start text-left font-normal",
                          !field.value && "text-muted-foreground"
                        )}
                      >
                        <CalendarIcon className="mr-2 h-4 w-4" />
                        {field.value ? (
                          format(field.value, "PPP")
                        ) : (
                          <span>Pick a date</span>
                        )}
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="flex w-auto flex-col space-y-2 p-2" align="end" side="left" sideOffset={20}>
                    <div className="flex justify-between">
                      <Select
                        value={calendarDate.getMonth().toString()}
                        onValueChange={handleMonthChange}
                      >
                        <SelectTrigger className="w-[140px]">
                          <SelectValue placeholder="Month" />
                        </SelectTrigger>
                        <SelectContent position="popper">
                          {months.map((month) => (
                            <SelectItem key={month.value} value={month.value}>
                              {month.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select
                        value={calendarDate.getFullYear().toString()}
                        onValueChange={handleYearChange}
                      >
                        <SelectTrigger className="w-[100px]">
                          <SelectValue placeholder="Year" />
                        </SelectTrigger>
                        <SelectContent position="popper">
                          {years.map((year) => (
                            <SelectItem key={year.value} value={year.value}>
                              {year.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="rounded-md border">
                        <Calendar
                          mode="single"
                          selected={field.value}
                          onSelect={field.onChange}
                          month={calendarDate}
                          onMonthChange={setCalendarDate}
                          initialFocus
                        />
                    </div>
                  </PopoverContent>
                </Popover>
                <FormDescription>
                  Select the date for your event.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <Button type="submit">Submit</Button>
        </form>
      </Form>
    )
}

