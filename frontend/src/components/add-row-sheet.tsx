import React from 'react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"

export default function AddRowSheet() {
    // TODO: add props and custom types
    // make the types have cells that return the type of component to use
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button className="ml-2">Add row</Button>
      </SheetTrigger>
      <SheetContent className="w-[400px] sm:w-[540px]">
        <SheetHeader>
          <SheetTitle>Edit profile</SheetTitle>
          <SheetDescription>
            Add brother information below then click "Done" to create new record.
          </SheetDescription>
        </SheetHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="firstName" className="text-right">
                First Name
            </Label>
            <Input id="firstName" className="col-span-3" />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="lastName" className="text-right">
                Last Name
            </Label>
            <Input id="lastName" className="col-span-3" />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="class" className="text-right">
                Status
            </Label>
            <Input id="class" className="col-span-3" />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="email" className="text-right">
                Email     
            </Label>
            <Input id="email" className="col-span-3" />
          </div>
        </div>
        <SheetFooter>
          <SheetClose asChild>
            <Button type="submit">Done</Button>
          </SheetClose>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}

