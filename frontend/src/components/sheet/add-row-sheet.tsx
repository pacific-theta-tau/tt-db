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

import FieldComponent, { FieldComponentProps } from './fields'
import { BrotherForm } from './row-form'

export default function AddRowSheet({ fields }: FieldComponentProps<TProps>) {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button className="ml-2">Add row</Button>
      </SheetTrigger>
      <SheetContent className="w-[400px] sm:w-[540px]">
        <SheetHeader>
          <SheetTitle>Create new Brother record</SheetTitle>
          <SheetDescription>
            Add the information below then click "Submit" to create new record.
          </SheetDescription>
        </SheetHeader>
            <BrotherForm />
        <SheetFooter>
            {/**/}
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}
