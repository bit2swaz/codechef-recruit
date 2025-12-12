import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { deletePost } from "@/lib/api"
import { Post } from "@/types"

export function useDeletePost(userId: number) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: deletePost,
    onMutate: async (postId) => {
      await queryClient.cancelQueries({ queryKey: ["posts", userId] })

      const previousPosts = queryClient.getQueryData<Post[]>(["posts", userId])

      queryClient.setQueryData(["posts", userId], (old: Post[] | undefined) => {
        return old ? old.filter((p) => p.id !== postId) : []
      })

      return { previousPosts }
    },
    onError: (err, newPost, context) => {
      if (context?.previousPosts) {
        queryClient.setQueryData(["posts", userId], context.previousPosts)
      }
      toast.error("Failed to delete post")
    },
    onSuccess: () => {
      toast.success("Post deleted successfully")
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["posts", userId] })
    },
  })
}
