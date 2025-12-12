"use client"

import { useState } from "react"
import { User, Bell, Shield, AlertTriangle } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { cn } from "@/lib/utils"

type Tab = "general" | "notifications" | "platform"

const CustomSwitch = ({ 
  checked, 
  onCheckedChange, 
  label, 
  description 
}: { 
  checked: boolean, 
  onCheckedChange: () => void, 
  label: string, 
  description?: string 
}) => (
  <div className="flex items-center justify-between gap-4 py-4">
    <div className="space-y-0.5 flex-1 min-w-0">
      <div className="font-medium text-base truncate">{label}</div>
      {description && <div className="text-sm text-muted-foreground break-words">{description}</div>}
    </div>
    <button
      onClick={onCheckedChange}
      className={cn(
        "relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
        checked ? "bg-primary" : "bg-input"
      )}
    >
      <span
        className={cn(
          "pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform duration-200 ease-in-out",
          checked ? "translate-x-5" : "translate-x-0"
        )}
      />
    </button>
  </div>
)

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState<Tab>("general")
  const [isLoading, setIsLoading] = useState(false)

  const [name, setName] = useState("Admin User")
  const [email, setEmail] = useState("admin@codechef-recruit.com")
  const [bio, setBio] = useState("Managing the community, one post at a time.")
  
  const [toggles, setToggles] = useState({
    emailAlerts: true,
    pushNotifs: false,
    weeklyDigest: true,
    maintenanceMode: false,
    registrations: true,
    autoApprove: false,
    strictFilter: true
  })

  const handleSave = () => {
    setIsLoading(true)
    setTimeout(() => {
      setIsLoading(false)
      toast.success("Settings saved successfully", {
        description: "Your changes have been applied to the system."
      })
    }, 1000)
  }

  const toggleSwitch = (key: keyof typeof toggles) => {
    setToggles(prev => {
      const newState = { ...prev, [key]: !prev[key] }
      
      if (key === 'maintenanceMode' && newState.maintenanceMode) {
        toast.warning("Maintenance Mode Enabled", {
          description: "The platform is now inaccessible to regular users."
        })
      }
      
      return newState
    })
  }

  return (
    <div className="h-full overflow-y-auto p-4 md:p-6 space-y-6 md:space-y-8 scroll-smooth pb-20">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-bold tracking-tight">Settings</h2>
          <p className="text-sm md:text-base text-muted-foreground">
            Manage your profile and platform preferences.
          </p>
        </div>
        <Button onClick={handleSave} disabled={isLoading} className="hidden md:flex">
          {isLoading ? "Saving..." : "Save Changes"}
        </Button>
      </div>

      <div className="grid grid-cols-3 gap-1 rounded-xl bg-muted/50 p-1 backdrop-blur-sm">
        {[
          { id: "general", label: "General", mobileLabel: "General", icon: User },
          { id: "notifications", label: "Notifications", mobileLabel: "Alerts", icon: Bell },
          { id: "platform", label: "Platform Controls", mobileLabel: "Platform", icon: Shield },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as Tab)}
            className={cn(
              "flex items-center justify-center gap-2 w-full rounded-lg px-2 py-2 md:px-3 md:py-2.5 text-xs md:text-sm font-medium leading-5 ring-white/60 ring-offset-2 ring-offset-blue-400 focus:outline-none focus:ring-2 transition-all duration-200",
              activeTab === tab.id
                ? "bg-background text-foreground shadow-sm scale-[1.02]"
                : "text-muted-foreground hover:bg-background/50 hover:text-foreground"
            )}
          >
            <tab.icon className="h-3 w-3 md:h-4 md:w-4" />
            <span className="hidden md:inline">{tab.label}</span>
            <span className="md:hidden">{tab.mobileLabel}</span>
          </button>
        ))}
      </div>

      <div className="grid gap-6">
        {activeTab === "general" && (
          <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <Card className="shadow-sm">
              <CardHeader>
                <CardTitle>Profile Information</CardTitle>
                <CardDescription>Update your public profile details.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="flex flex-col md:flex-row gap-6 items-start">
                  <div className="flex flex-col items-center gap-3">
                    <Avatar className="h-24 w-24 border-4 border-background shadow-xl">
                      <AvatarImage src="/placeholder-avatar.jpg" />
                      <AvatarFallback className="text-2xl bg-primary/10 text-primary">AD</AvatarFallback>
                    </Avatar>
                    <Button variant="outline" size="sm" className="w-full">Change Avatar</Button>
                  </div>
                  <div className="flex-1 space-y-4 w-full">
                    <div className="grid gap-2">
                      <label className="text-sm font-medium">Display Name</label>
                      <Input 
                        value={name} 
                        onChange={(e) => setName(e.target.value)} 
                        className="bg-background/50"
                      />
                    </div>
                    <div className="grid gap-2">
                      <label className="text-sm font-medium">Email</label>
                      <Input 
                        value={email} 
                        onChange={(e) => setEmail(e.target.value)} 
                        className="bg-background/50"
                      />
                    </div>
                    <div className="grid gap-2">
                      <label className="text-sm font-medium">Bio</label>
                      <textarea 
                        className="flex min-h-[80px] w-full rounded-md border border-input bg-background/50 px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        value={bio}
                        onChange={(e) => setBio(e.target.value)}
                      />
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {activeTab === "notifications" && (
          <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <Card className="shadow-sm">
              <CardHeader>
                <CardTitle>Alert Preferences</CardTitle>
                <CardDescription>Choose how you want to be notified about platform activity.</CardDescription>
              </CardHeader>
              <CardContent className="divide-y divide-border/40">
                <CustomSwitch
                  label="Email Alerts"
                  description="Receive emails about critical system events."
                  checked={toggles.emailAlerts}
                  onCheckedChange={() => toggleSwitch('emailAlerts')}
                />
                <CustomSwitch
                  label="Push Notifications"
                  description="Receive real-time push notifications on your device."
                  checked={toggles.pushNotifs}
                  onCheckedChange={() => toggleSwitch('pushNotifs')}
                />
                <CustomSwitch
                  label="Weekly Digest"
                  description="Get a weekly summary of platform statistics."
                  checked={toggles.weeklyDigest}
                  onCheckedChange={() => toggleSwitch('weeklyDigest')}
                />
              </CardContent>
            </Card>
          </div>
        )}

        {activeTab === "platform" && (
          <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
            <Card className="shadow-sm border-l-4 border-l-primary">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <Shield className="h-5 w-5 text-primary" />
                  <CardTitle>Administration Controls</CardTitle>
                </div>
                <CardDescription>Manage global platform settings. Use with caution.</CardDescription>
              </CardHeader>
              <CardContent className="divide-y divide-border/40">
                <CustomSwitch
                  label="Maintenance Mode"
                  description="Disable access for all non-admin users."
                  checked={toggles.maintenanceMode}
                  onCheckedChange={() => toggleSwitch('maintenanceMode')}
                />
                <CustomSwitch
                  label="Allow New Registrations"
                  description="Open or close the platform for new users."
                  checked={toggles.registrations}
                  onCheckedChange={() => toggleSwitch('registrations')}
                />
                <CustomSwitch
                  label="Auto-approve Posts"
                  description="Skip manual moderation for new content."
                  checked={toggles.autoApprove}
                  onCheckedChange={() => toggleSwitch('autoApprove')}
                />
                <CustomSwitch
                  label="Strict Content Filter"
                  description="Enable AI-powered content moderation."
                  checked={toggles.strictFilter}
                  onCheckedChange={() => toggleSwitch('strictFilter')}
                />
              </CardContent>
            </Card>

            <Card className="bg-destructive/5 backdrop-blur-xl border-destructive/20 shadow-sm">
              <CardHeader>
                <div className="flex items-center gap-2 text-destructive">
                  <AlertTriangle className="h-5 w-5" />
                  <CardTitle>Danger Zone</CardTitle>
                </div>
                <CardDescription>Irreversible actions for the platform.</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between py-4">
                  <div className="space-y-0.5">
                    <div className="font-medium text-base">Purge All Deleted Data</div>
                    <div className="text-sm text-muted-foreground">Permanently remove soft-deleted posts and users.</div>
                  </div>
                  <Button variant="destructive" onClick={() => toast.error("Action blocked", { description: "You do not have permission to perform this action." })}>
                    Purge Data
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>

      <div className="md:hidden">
        <Button onClick={handleSave} disabled={isLoading} className="w-full">
          {isLoading ? "Saving..." : "Save Changes"}
        </Button>
      </div>
    </div>
  )
}
