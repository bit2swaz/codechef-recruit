"use client"

import { Menu } from "lucide-react"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { Sidebar } from "@/components/sidebar"
import { Button } from "@/components/ui/button"
import { useState, useEffect } from "react"
import { usePathname } from "next/navigation"

export function MobileSidebar() {
  const [open, setOpen] = useState(false)
  const pathname = usePathname()

  // Close sidebar when route changes
  useEffect(() => {
    setOpen(false)
  }, [pathname])

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant="ghost" size="icon" className="md:hidden mr-2">
          <Menu />
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="p-0 w-72">
        <Sidebar className="w-full border-r-0" />
      </SheetContent>
    </Sheet>
  )
}
