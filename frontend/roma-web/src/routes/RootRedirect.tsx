import useAuth from "@/store/auth"
import { Navigate } from "react-router-dom"

export default function RootRedirect() { 
      const { isAuthenticated } = useAuth()
  return <Navigate to={isAuthenticated ? '/sessions' : '/auth/login'} replace />
 }