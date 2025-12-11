"use client"

import { Bell, Search, User } from "lucide-react"
import Link from "next/link"
import { toast } from "sonner"
import { MobileSidebar } from "@/components/mobile-sidebar"
import { ModeToggle } from "@/components/mode-toggle"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export function Header() {
  return (
    <header className="h-16 border-b border-border/40 flex items-center justify-between px-6 bg-background/60 backdrop-blur-xl sticky top-0 z-50 transition-all duration-300">
      <div className="flex items-center gap-4 lg:gap-8">
        <MobileSidebar />
        <div className="hidden md:flex items-center gap-2 text-muted-foreground">
          <span className="font-semibold text-foreground">Dashboard</span>
          <span className="text-muted-foreground/50">/</span>
          <span className="text-sm">Overview</span>
        </div>
      </div>
      
      <div className="flex items-center gap-4">
        <div className="relative hidden md:block w-64">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input 
            placeholder="Global search..." 
            className="pl-9 bg-background/50 border-border/50 focus-visible:ring-primary/20 transition-all duration-300" 
          />
        </div>

        <div className="flex items-center gap-2">
          <Button 
            variant="ghost" 
            size="icon" 
            className="relative hover:bg-primary/10 transition-colors"
            onClick={() => toast.info("This button is a placeholder and doesn't do anything yet.")}
          >
            <Bell className="h-5 w-5" />
            <span className="absolute top-2 right-2 h-2 w-2 bg-primary rounded-full animate-pulse" />
          </Button>
          <ModeToggle />
          
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="relative h-9 w-9 rounded-full ml-2 ring-2 ring-transparent hover:ring-primary/20 transition-all">
                <Avatar className="h-9 w-9 border border-border/50">
                  <AvatarImage src="/placeholder-avatar.jpg" alt="@admin" />
                  <AvatarFallback className="bg-primary/10 text-primary font-medium">AD</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium leading-none">Admin User</p>
                  <p className="text-xs leading-none text-muted-foreground">
                    admin@codechef-recruit.com
                  </p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="cursor-pointer" onClick={() => toast.info("This button is a placeholder and doesn't do anything yet.")}>
                <span>Profile</span>
              </DropdownMenuItem>
              <DropdownMenuItem asChild className="cursor-pointer">
                <Link href="/dashboard/settings">Settings</Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="text-destructive focus:text-destructive cursor-pointer">
                Log out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}
