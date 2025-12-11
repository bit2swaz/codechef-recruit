import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { deletePost } from "@/lib/api"
import { Post } from "@/types"

export function useDeletePost(userId: number) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: deletePost,
    onMutate: async (postId) => {
      // Cancel any outgoing refetches (so they don't overwrite our optimistic update)
      await queryClient.cancelQueries({ queryKey: ["posts", userId] })

      // Snapshot the previous value
      const previousPosts = queryClient.getQueryData<Post[]>(["posts", userId])

      // Optimistically update to the new value
      queryClient.setQueryData(["posts", userId], (old: Post[] | undefined) => {
        return old ? old.filter((p) => p.id !== postId) : []
      })

      // Return a context object with the snapshotted value
      return { previousPosts }
    },
    onError: (err, newPost, context) => {
      // If the mutation fails, use the context returned from onMutate to roll back
      if (context?.previousPosts) {
        queryClient.setQueryData(["posts", userId], context.previousPosts)
      }
      toast.error("Failed to delete post")
    },
    onSuccess: () => {
      toast.success("Post deleted successfully")
    },
    onSettled: () => {
      // Always refetch after error or success:
      queryClient.invalidateQueries({ queryKey: ["posts", userId] })
    },
  })
}
