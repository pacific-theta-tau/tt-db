import React, { useState } from 'react'
import { Trash2 } from "lucide-react"
import { Button } from "@/components/ui/button"
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
}

export const DeleteAlertDialog: React.FC<DeleteAlertDialogProps> = ({ trigger, endpoint, body }) => {
    const [isOpen, setIsOpen] = useState(false)
    const [isLoading, setIsLoading] = useState(false);
    const { toast } = useToast()
    const handleDelete = async () => {
        setIsLoading(true)
        console.log("HANDLE DELETE FUNCTION CALLED")
        try {
            request(endpoint, "DELETE", body)
            await new Promise( resolve => setTimeout(resolve, 5000) );
            toast({
                title: "Deleted Successfully",
                description: "The row has been deleted successfully.",
            })
        } catch (error) {
            toast({
                title: "Failed to delete row",
                description: `Failed to delete the item. Please try again.\nError: ${error}`,
                variant: "destructive",
            })
        } finally {
            setIsLoading(false)
            setIsOpen(false)
        }
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
                  <AlertDialogAction onClick={handleDelete} disabled={isLoading}>
                    { isLoading ? "Deleting..." : "Delete" }
                  </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )

}
