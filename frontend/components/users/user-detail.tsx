"use client"

import { useQuery } from "@tanstack/react-query"
import { Mail, Globe, Building, ArrowLeft, MessageSquare, Trash2 } from "lucide-react"
import { useRouter } from "next/navigation"

import { getPosts } from "@/lib/api"
import { useDeletePost } from "@/hooks/use-delete-post"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { User } from "@/types"

interface UserDetailProps {
  selectedUserId: string | null
  user: User | undefined
}

export function UserDetail({ selectedUserId, user }: UserDetailProps) {
  const router = useRouter()
  const { mutate: deletePost } = useDeletePost(Number(selectedUserId))

  const { data: posts, isLoading: isLoadingPosts } = useQuery({
    queryKey: ["posts", Number(selectedUserId)],
    queryFn: () => getPosts(Number(selectedUserId)),
    enabled: !!selectedUserId,
  })

  if (!selectedUserId) {
    return (
      <div className="flex h-full flex-col items-center justify-center text-muted-foreground p-8 text-center bg-background/60 backdrop-blur-xl">
        <div className="h-24 w-24 rounded-full bg-muted/30 flex items-center justify-center mb-6 animate-pulse">
          <MessageSquare className="h-10 w-10 opacity-50" />
        </div>
        <h3 className="text-xl font-semibold mb-2 text-foreground">No User Selected</h3>
        <p className="text-sm max-w-xs text-muted-foreground/80">
          Select a user from the list to view their profile and recent posts.
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full bg-background/60 backdrop-blur-xl">
      {user && (
        <div className="p-6 border-b border-border/40 bg-background/80 backdrop-blur-md z-10">
          <div className="flex items-start gap-4">
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden -ml-2 hover:bg-background/80"
              onClick={() => router.push("/dashboard/users")}
            >
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div className="flex-1 flex items-start justify-between">
              <div className="space-y-1.5">
                <h2 className="text-2xl font-bold tracking-tight text-foreground">
                  {user.name}
                </h2>
                <div className="flex flex-col gap-1 sm:flex-row sm:gap-4 text-sm text-muted-foreground">
                  <div className="flex items-center gap-2 hover:text-primary transition-colors">
                    <Mail className="h-3.5 w-3.5" />
                    {user.email}
                  </div>
                  <div className="flex items-center gap-2 hover:text-primary transition-colors">
                    <Globe className="h-3.5 w-3.5" />
                    {user.website}
                  </div>
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground pt-1 hover:text-primary transition-colors">
                  <Building className="h-3.5 w-3.5" />
                  {user.company.name}
                </div>
              </div>
              <Avatar className="h-16 w-16 border-2 border-border/50 shadow-lg hidden sm:block transition-transform hover:scale-105 duration-300">
                <AvatarFallback className="text-xl bg-primary/10 text-primary font-bold">
                  {user.name.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
            </div>
          </div>
        </div>
      )}

      <ScrollArea className="h-full">
        <div className="p-6 max-w-3xl mx-auto space-y-6">
          <h3 className="text-lg font-semibold flex items-center gap-2 text-foreground/80">
            <MessageSquare className="h-4 w-4 text-primary" />
            Recent Posts
          </h3>
          {isLoadingPosts ? (
            <div className="space-y-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <Card key={i} className="border-border/40 shadow-none bg-card/40 backdrop-blur-sm">
                  <CardHeader>
                    <Skeleton className="h-6 w-3/4" />
                  </CardHeader>
                  <CardContent>
                    <Skeleton className="h-4 w-full mb-2" />
                    <Skeleton className="h-4 w-5/6" />
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <div className="grid gap-4">
              {posts?.map((post) => (
                <Card 
                  key={post.id} 
                  className="group hover:shadow-lg transition-all duration-300 border-border/40 bg-card/40 backdrop-blur-sm hover:bg-card/60 hover:scale-[1.01]"
                >
                  <CardHeader className="pr-12">
                    <CardTitle className="text-base font-semibold group-hover:text-primary transition-colors duration-300">
                      {post.title}
                    </CardTitle>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="absolute top-4 right-4 text-muted-foreground hover:text-destructive hover:bg-destructive/10 opacity-0 group-hover:opacity-100 transition-all duration-200 scale-90 group-hover:scale-100"
                      onClick={() => deletePost(post.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground leading-relaxed group-hover:text-foreground/80 transition-colors">
                      {post.body}
                    </p>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </ScrollArea>
    </div>
  )
}
