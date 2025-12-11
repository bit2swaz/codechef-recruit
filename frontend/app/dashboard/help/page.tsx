import { Construction } from "lucide-react"

export default function HelpPage() {
  return (
    <div className="flex h-full flex-col items-center justify-center text-center p-8 space-y-4">
      <div className="h-24 w-24 rounded-full bg-muted/30 flex items-center justify-center animate-pulse">
        <Construction className="h-12 w-12 text-muted-foreground" />
      </div>
      <h2 className="text-2xl font-bold tracking-tight">Coming Soon</h2>
      <p className="text-muted-foreground max-w-md">
        We are currently building a comprehensive help center to assist you with managing the platform. Check back later!
      </p>
    </div>
  )
}
