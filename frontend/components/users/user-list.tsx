"use client"

import { Virtuoso } from "react-virtuoso"
import { Search } from "lucide-react"
import { useRouter } from "next/navigation"

import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Skeleton } from "@/components/ui/skeleton"
import { cn } from "@/lib/utils"
import { User } from "@/types"

interface UserListProps {
  users: User[] | undefined
  isLoading: boolean
  search: string
  setSearch: (value: string) => void
  selectedUserId: string | null
}

export function UserList({
  users,
  isLoading,
  search,
  setSearch,
  selectedUserId,
}: UserListProps) {
  const router = useRouter()

  const filteredUsers = users?.filter((user) =>
    user.name.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="flex flex-col h-full bg-background/60 backdrop-blur-xl">
      <div className="p-4 border-b border-border/40 sticky top-0 bg-background/80 backdrop-blur-md z-10">
        <div className="relative group">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
          <Input
            placeholder="Search users..."
            className="pl-9 bg-background/50 border-border/50 focus-visible:ring-primary/20 focus-visible:border-primary/50 transition-all duration-300"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>
      <div className="flex-1 h-full">
        {isLoading ? (
          <div className="p-4 space-y-4">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3 animate-pulse">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="space-y-2 flex-1">
                  <Skeleton className="h-4 w-[140px]" />
                  <Skeleton className="h-3 w-[100px]" />
                </div>
              </div>
            ))}
          </div>
        ) : (
          <Virtuoso
            data={filteredUsers}
            itemContent={(index, user) => (
              <div
                key={user.id}
                className={cn(
                  "flex items-center gap-3 p-4 border-b border-border/40 cursor-pointer transition-all duration-200 hover:pl-5",
                  selectedUserId === user.id.toString() 
                    ? "bg-primary/10 border-l-4 border-l-primary pl-5" 
                    : "hover:bg-muted/30 hover:border-l-4 hover:border-l-muted-foreground/30"
                )}
                onClick={() => router.push(`/dashboard/users?userId=${user.id}`)}
              >
                <Avatar className={cn(
                  "h-10 w-10 border transition-transform duration-300",
                  selectedUserId === user.id.toString() ? "scale-110 border-primary/50" : "border-border/50"
                )}>
                  <AvatarFallback className="bg-primary/10 text-primary font-medium">
                    {user.name.slice(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col overflow-hidden">
                  <span className={cn(
                    "font-medium truncate text-sm transition-colors",
                    selectedUserId === user.id.toString() ? "text-primary" : "text-foreground"
                  )}>{user.name}</span>
                  <span className="text-xs text-muted-foreground truncate">
                    {user.email}
                  </span>
                </div>
              </div>
            )}
          />
        )}
      </div>
    </div>
  )
}
