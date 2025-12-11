import { Sidebar } from "@/components/sidebar"
import { Header } from "@/components/header"

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="flex h-screen overflow-hidden bg-background relative">
      <Sidebar className="hidden md:flex" />
      <div className="flex-1 flex flex-col relative z-0">
        <Header />
        <main className="flex-1 overflow-hidden flex flex-col">
          {children}
        </main>
      </div>
    </div>
  )
}
