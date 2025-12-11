"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { LayoutDashboard, Users, Settings, LogOut, HelpCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

const routes = [
  {
    label: "Dashboard",
    icon: LayoutDashboard,
    href: "/dashboard",
  },
  {
    label: "Users",
    icon: Users,
    href: "/dashboard/users",
  },
  {
    label: "Settings",
    icon: Settings,
    href: "/dashboard/settings",
  },
]

export function Sidebar({ className }: { className?: string }) {
  const pathname = usePathname()

  return (
    <aside className={cn(
      "w-64 h-full border-r flex flex-col bg-background/60 backdrop-blur-xl border-border/40 transition-all duration-300",
      className
    )}>
      <div className="h-16 flex items-center px-6 border-b border-border/40">
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/20 flex items-center justify-center">
            <div className="h-4 w-4 rounded-sm bg-primary animate-pulse" />
          </div>
          <h1 className="text-xl font-bold text-foreground">
            AdminPanel
          </h1>
        </div>
      </div>
      
      <div className="flex-1 py-6 px-3 space-y-1 overflow-y-auto">
        <div className="px-3 mb-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Main Menu
        </div>
        {routes.map((route) => (
          <Link key={route.href} href={route.href}>
            <Button
              variant="ghost"
              className={cn(
                "w-full justify-start mb-1 transition-all duration-200 hover:scale-[1.02] hover:bg-primary/10",
                pathname === route.href 
                  ? "bg-primary/15 text-primary font-medium shadow-sm border-r-2 border-primary rounded-r-none" 
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              <route.icon className={cn(
                "mr-3 h-4 w-4 transition-colors",
                pathname === route.href ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
              )} />
              {route.label}
            </Button>
          </Link>
        ))}

        <div className="px-3 mt-8 mb-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Support
        </div>
        <Link href="/dashboard/help">
          <Button variant="ghost" className={cn(
            "w-full justify-start mb-1 transition-all duration-200 hover:scale-[1.02] hover:bg-primary/10",
            pathname === "/dashboard/help"
              ? "bg-primary/15 text-primary font-medium shadow-sm border-r-2 border-primary rounded-r-none"
              : "text-muted-foreground hover:text-foreground"
          )}>
            <HelpCircle className={cn(
              "mr-3 h-4 w-4 transition-colors",
              pathname === "/dashboard/help" ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
            )} />
            Help Center
          </Button>
        </Link>
      </div>

      <div className="p-4 border-t border-border/40 bg-background/40">
        <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-background/60 transition-colors cursor-pointer group">
          <Avatar className="h-9 w-9 border border-border/50 group-hover:border-primary/50 transition-colors">
            <AvatarImage src="/placeholder-avatar.jpg" />
            <AvatarFallback className="bg-primary/10 text-primary">AD</AvatarFallback>
          </Avatar>
          <div className="flex-1 overflow-hidden">
            <p className="text-sm font-medium truncate group-hover:text-primary transition-colors">Admin User</p>
            <p className="text-xs text-muted-foreground truncate">admin@codechef.com</p>
          </div>
          <LogOut className="h-4 w-4 text-muted-foreground hover:text-destructive transition-colors" />
        </div>
      </div>
    </aside>
  )
}
