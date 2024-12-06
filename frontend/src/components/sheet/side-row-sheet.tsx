// add-row-sheet.tsx: Sheet+form that opens when "Add row" button is clicked in data table pages
import React from 'react'
import { Button } from "@/components/ui/button"
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


const SideFormSheet: React.FC<{
    title: string;
    description: string;
    FormType: React.JSX.Element
    trigger?: React.ReactNode
}> = ({ title, description, FormType, trigger }) => {
  return (
    <Sheet>
      <SheetTrigger asChild>
        { trigger ?
            trigger : 
            <Button className="ml-2">Add Row</Button>
        }
      </SheetTrigger>
          <SheetContent className="w-[400px] sm:w-[540px] overflow-y-auto">
                <SheetHeader>
                  <SheetTitle>{title}</SheetTitle>
                  <SheetDescription>
                    {description}
                  </SheetDescription>
                </SheetHeader>
                    { FormType }
                <SheetFooter>
                    {/**/}
                </SheetFooter>
          </SheetContent>
    </Sheet>
  )
}

export default SideFormSheet 
