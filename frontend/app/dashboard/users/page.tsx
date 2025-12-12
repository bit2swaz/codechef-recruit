"use client"

import { useState, Suspense } from "react"
import { useQuery } from "@tanstack/react-query"
import { useSearchParams } from "next/navigation"

import { getUsers } from "@/lib/api"
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import { UserList } from "@/components/users/user-list"
import { UserDetail } from "@/components/users/user-detail"

function UsersPageContent() {
  const searchParams = useSearchParams()
  const selectedUserId = searchParams.get("userId")
  const [search, setSearch] = useState("")

  const { data: users, isLoading: isLoadingUsers } = useQuery({
    queryKey: ["users"],
    queryFn: getUsers,
  })

  const selectedUser = users?.find((u) => u.id.toString() === selectedUserId)

  return (
    <div className="h-full bg-background flex flex-col overflow-hidden">
      <div className="md:hidden h-full">
        {!selectedUserId ? (
          <UserList
            users={users}
            isLoading={isLoadingUsers}
            search={search}
            setSearch={setSearch}
            selectedUserId={selectedUserId}
          />
        ) : (
          <UserDetail selectedUserId={selectedUserId} user={selectedUser} />
        )}
      </div>

      <div className="hidden md:block h-full">
        <ResizablePanelGroup direction="horizontal" className="h-full items-stretch">
          <ResizablePanel
            defaultSize={35}
            minSize={25}
            maxSize={45}
            className="border-r"
          >
            <UserList
              users={users}
              isLoading={isLoadingUsers}
              search={search}
              setSearch={setSearch}
              selectedUserId={selectedUserId}
            />
          </ResizablePanel>
          
          <ResizableHandle withHandle />
          
          <ResizablePanel defaultSize={65}>
            <UserDetail selectedUserId={selectedUserId} user={selectedUser} />
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  )
}

export default function UsersPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center h-full">Loading...</div>}>
      <UsersPageContent />
    </Suspense>
  )
}
