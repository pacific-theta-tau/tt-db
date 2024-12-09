import React, { useState } from 'react'
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { request } from '@/api/api';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { useToast } from "@/hooks/use-toast"


interface DeleteAlertDialogProps {
    trigger?: React.ReactNode
    endpoint: string
    body?: Record<string, unknown>
    queryKey?: string
}

export const DeleteAlertDialog: React.FC<DeleteAlertDialogProps> = ({ trigger, endpoint, body, queryKey }) => {
    // Used to keep Dialog component open until request is done
    const [isOpen, setIsOpen] = useState(false)
    const { toast } = useToast()

    // React Query and Mutation hooks
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationFn: async () => {
            const data = await request(endpoint, "DELETE", body)
            console.log(data)
        },
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: [queryKey] });
            toast({
                title: "Deleted Successfully",
                description: "The row has been deleted successfully.",
            })
        },
        onError: (error) => {
            console.error('Error fetching data:', error);
            toast({
                title: "Failed to delete row",
                description: error instanceof Error ? error.message : "Failed to delete the item. Please try again.",
                variant: "destructive",
            })

        },
        onSettled: () => {
            setIsOpen(false)
        }
    })

    const handleDelete = async () => {
        mutation.mutate()
    }

    return (
        <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
            <AlertDialogTrigger asChild>
                { trigger }
            </AlertDialogTrigger>
            <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This action cannot be undone. This will permanently delete the row from the Database
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction onClick={handleDelete} disabled={mutation.isPending}>
                    { mutation.isPending ? "Deleting..." : "Delete" }
                  </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )

}
