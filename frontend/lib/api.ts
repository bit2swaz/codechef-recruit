import { Post, User } from "@/types"

const BASE_URL = "https://jsonplaceholder.typicode.com"

async function fetcher<T>(url: string): Promise<T> {
  await new Promise((resolve) => setTimeout(resolve, 800))

  const response = await fetch(`${BASE_URL}${url}`)

  if (!response.ok) {
    throw new Error(`API Error: ${response.statusText}`)
  }

  return response.json()
}

export async function getUsers(): Promise<User[]> {
  return fetcher<User[]>("/users")
}

export async function getPosts(userId: number): Promise<Post[]> {
  return fetcher<Post[]>(`/posts?userId=${userId}`)
}

export async function deletePost(postId: number): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 800))
  
  const response = await fetch(`${BASE_URL}/posts/${postId}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error(`API Error: ${response.statusText}`)
  }
}
