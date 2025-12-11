"use client"

import { useQuery } from "@tanstack/react-query"
import { Activity, CreditCard, Server, Users, AlertTriangle, ShieldAlert, UserPlus, MessageSquare, Zap } from "lucide-react"
import { getUsers } from "@/lib/api"
import { StatsCard } from "@/components/dashboard/stats-card"
import { OverviewChart } from "@/components/dashboard/overview-chart"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { ScrollArea } from "@/components/ui/scroll-area"
import { toast } from "sonner"

export default function DashboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["users"],
    queryFn: getUsers,
  })

  const recentActivity = [
    { user: "Alice Cooper", action: "reported a post", time: "2 mins ago", type: "warning" },
    { user: "Bob Smith", action: "registered new account", time: "5 mins ago", type: "success" },
    { user: "Charlie Brown", action: "updated profile picture", time: "12 mins ago", type: "neutral" },
    { user: "System", action: "Automated backup completed", time: "1 hour ago", type: "system" },
    { user: "David Lee", action: "posted in 'General'", time: "2 hours ago", type: "neutral" },
  ]

  return (
    <div className="h-full overflow-y-auto p-6 space-y-6 scroll-smooth pb-20">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-foreground">
            Command Center
          </h2>
          <p className="text-muted-foreground">
            Real-time overview of platform activity and health.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button 
            variant="outline" 
            className="gap-2 bg-background/50 backdrop-blur-sm h-10 px-3 md:px-4"
            onClick={() => toast.info("This button is a placeholder and doesn't do anything yet.")}
          >
            <Zap className="h-4 w-4 text-primary" />
            <span className="hidden md:inline">Quick Actions</span>
          </Button>
          <Button 
            className="gap-2 shadow-lg shadow-primary/20 h-10 px-3 md:px-4"
            onClick={() => toast.info("This button is a placeholder and doesn't do anything yet.")}
          >
            <ShieldAlert className="h-4 w-4" />
            <span className="hidden md:inline">Moderation Queue</span>
            <Badge variant="secondary" className="ml-1 bg-primary-foreground text-primary hover:bg-primary-foreground/90">12</Badge>
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-[120px] rounded-xl" />
          ))
        ) : (
          <>
            <StatsCard
              title="Total Users"
              value={data?.length || 0}
              icon={Users}
              description="+20.1% from last month"
            />
            <StatsCard
              title="Active Sessions"
              value="1,203"
              icon={Activity}
              description="+201 since last hour"
            />
            <StatsCard
              title="Pending Reports"
              value="12"
              icon={AlertTriangle}
              description="Requires immediate attention"
            />
            <StatsCard
              title="Server Health"
              value="98.2%"
              icon={Server}
              description="All systems operational"
            />
          </>
        )}
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4 shadow-sm hover:shadow-lg hover:scale-[1.005] transition-all duration-300">
          <CardHeader>
            <CardTitle>User Growth & Engagement</CardTitle>
            <CardDescription>Daily active users vs new registrations over time.</CardDescription>
          </CardHeader>
          <CardContent className="pl-2">
            <OverviewChart />
          </CardContent>
        </Card>

        <div className="col-span-3 space-y-4">
          <Card className="h-full shadow-sm hover:shadow-lg hover:scale-[1.005] transition-all duration-300 flex flex-col">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Activity className="h-5 w-5 text-primary" />
                Live Activity Feed
              </CardTitle>
              <CardDescription>Real-time actions across the platform.</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 p-0">
              <ScrollArea className="h-[300px] px-6">
                <div className="space-y-6 pb-6">
                  {recentActivity.map((item, i) => (
                    <div key={i} className="flex items-start gap-4 group">
                      <div className="relative mt-1">
                        <span className="absolute left-2 top-8 h-full w-[1px] bg-border group-last:hidden" />
                        <div className={`h-4 w-4 rounded-full border-2 border-background ring-2 ${
                          item.type === 'warning' ? 'bg-destructive ring-destructive/20' :
                          item.type === 'success' ? 'bg-green-500 ring-green-500/20' :
                          item.type === 'system' ? 'bg-blue-500 ring-blue-500/20' :
                          'bg-muted-foreground ring-muted/20'
                        }`} />
                      </div>
                      <div className="space-y-1">
                        <p className="text-sm font-medium leading-none">
                          <span className="text-primary">{item.user}</span> {item.action}
                        </p>
                        <p className="text-xs text-muted-foreground">{item.time}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </ScrollArea>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card className="hover:shadow-lg hover:scale-[1.02] transition-all duration-300 cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-foreground">
              <UserPlus className="h-5 w-5" />
              New Signups
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">+128</div>
            <p className="text-xs text-muted-foreground">Today</p>
          </CardContent>
        </Card>
        <Card className="hover:shadow-lg hover:scale-[1.02] transition-all duration-300 cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-foreground">
              <MessageSquare className="h-5 w-5" />
              New Posts
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">+2,340</div>
            <p className="text-xs text-muted-foreground">Today</p>
          </CardContent>
        </Card>
        <Card className="hover:shadow-lg hover:scale-[1.02] transition-all duration-300 cursor-pointer">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-foreground">
              <Zap className="h-5 w-5" />
              API Latency
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">45ms</div>
            <p className="text-xs text-muted-foreground">Average</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
