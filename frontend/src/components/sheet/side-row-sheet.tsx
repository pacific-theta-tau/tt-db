// add-row-sheet.tsx: Sheet+form that opens when "Add row" button is clicked in data table pages
import React, { useState } from 'react'
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
  /* Controlling the open/close state of `Sheet` because some `FormType` components require
   * e.preventDefault(), which prevents the sheet from closing on form submission.
   * By passing the `setIsOpen` function to the FormType, we can triggeer closing the sheet after submission.
   */
  const [isOpen, setIsOpen] = useState(false);
  const closeDialog = () => setIsOpen(false);

  const additionalProps = { onClose: closeDialog };
  const newForm = React.cloneElement(FormType, additionalProps)

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
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
                    { /* FormType */ }
                    { newForm }
                <SheetFooter>
                    {/**/}
                </SheetFooter>
          </SheetContent>
    </Sheet>
  )
}

export default SideFormSheet 
